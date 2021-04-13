package trigger

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func StartRunFailsDiscovery(db repository.Database, tr utils.ConfigUtils, CommonMutex *utils.CommonMutex) {
	go discoveryRunFails(db, tr, CommonMutex)
}

func discoveryRunFails(db repository.Database, tr utils.ConfigUtils, CommonMutex *utils.CommonMutex) {
	for {
		log.Println("Start discoveryRunFails")

		organizationsName, _ := db.GetOrganizationsName()

		for _, organizationName := range organizationsName {
			mutex := utils.ReserveOrganizationMutex(organizationName, CommonMutex)
			mutex.Lock()

			locked := true
			var deferRun func(name string, voteMutex *utils.CommonMutex, mutex *sync.Mutex, locked *bool) = utils.ReleaseOrganizationMutexDefer
			defer deferRun(organizationName, CommonMutex, mutex, &locked)

			org, _ := db.GetOrganizationByName(organizationName)
			if org == nil {
				log.Println("discoveryRunFails organization ", organizationName, "not found")
				continue
			}

			gitSource, err := db.GetGitSourceById(org.GitSourceID)
			if gitSource == nil || err != nil || org.Projects == nil {
				log.Println("discoveryRunFails gitsource not fount for", org.Name, "organization")
				continue
			}

			for projectName, project := range org.Projects {
				if project.Archivied {
					continue
				}

				checkNewRuns := CheckIfNewRunsPresent(&project)
				if !checkNewRuns {
					log.Println("no new runs found for project", projectName)
					continue
				}

				//If there are new runs asks for other runs
				lastRun := project.GetLastRun()
				runList, _ := agola.GetRuns(project.AgolaProjectID, false, "finished", &lastRun.ID, 0, true)

				runList = takeWebhookTrigger(runList)

				for _, run := range *runList {
					newRun := model.RunInfo{
						ID:           run.ID,
						Branch:       run.GetBranchName(),
						RunStartDate: run.StartTime,
						RunEndDate:   run.EndTime,
						Phase:        model.RunPhase(run.Phase),
						Result:       model.RunResult(run.Result),
					}
					project.PushNewRun(newRun)

					//

					if run.Result == agola.RunResultFailed && run.StartTime.After(lastRun.RunStartDate) {
						log.Println("Found run failed!")
						emailMap := getUsersEmailMap(gitSource, org, project.GitRepoPath, run)
						log.Println("send emails to:", emailMap)

						body, err := makeBody(org.Name, project.GitRepoPath, run)
						if err != nil {
							log.Println("Failed to make email body")
							continue
						}
						subject := makeSubject(org.Name, project.GitRepoPath, run)

						sendConfirmEmail(emailMap, nil, subject, body)
					}
				}

				org.Projects[projectName] = project
			}
			db.SaveOrganization(org)

			mutex.Unlock()
			utils.ReleaseOrganizationMutex(organizationName, CommonMutex)
			locked = false
		}

		log.Println("End discoveryRunFails")
		time.Sleep(time.Duration(tr.GetRunFailedTriggerTime()) * time.Minute)
	}
}

func getUsersEmailMap(gitSource *model.GitSource, organization *model.Organization, gitRepoPath string, failedRun agola.RunDto) map[string]bool {
	emails := make(map[string]bool, 0)

	//Find all users that commited the failed run and parents
	emailUsersCommitted := getEmailByRun(&failedRun, gitSource, organization.Name, gitRepoPath)

	//Users owner of the organization and users owner of the repository
	var usersRepoOwners *[]string

	if gitSource.GitType == model.Gitea {
		usersRepoOwners, _ = findGiteaUsersEmailRepositoryOwner(gitSource, organization.Name, gitRepoPath)
	} else {
		usersRepoOwners, _ = findGithubUsersRepositoryOwner(gitSource, organization.Name, gitRepoPath)
	}

	for _, email := range emailUsersCommitted {
		emails[email] = true
	}

	if usersRepoOwners != nil {
		for _, email := range *usersRepoOwners {
			emails[email] = true
		}
	}

	if organization.ExternalUsers != nil {
		for email, _ := range organization.ExternalUsers {
			emails[email] = true
		}
	}

	return emails
}

const bodyMessageTemplate string = "[%s/%s] FIX Agola Run (#%s)\n"
const bodyLinkTemplate string = `See: <a href="%s">click here</a>`
const subjectTemplate string = "Run failed in Agola: %s » %s » release #%s"
const runAgolaPath string = "%s/org/%s/projects/%s.proj/runs/%s"

func makeSubject(organizationName string, projectName string, failedRun agola.RunDto) string {
	return fmt.Sprintf(subjectTemplate, organizationName, projectName, fmt.Sprint(failedRun.Counter))
}

func getRunAgolaUrl(organizationName string, projectName string, runID string) string {
	return fmt.Sprintf(runAgolaPath, config.Config.Agola.AgolaAddr, organizationName, projectName, runID)
}

func makeBody(organizationName string, projectName string, failedRun agola.RunDto) (string, error) {
	runUrl := getRunAgolaUrl(organizationName, projectName, failedRun.ID)
	body := fmt.Sprintf(bodyMessageTemplate, organizationName, projectName, fmt.Sprint(failedRun.Counter))
	body += fmt.Sprintf(bodyLinkTemplate, runUrl)

	run, err := agola.GetRun(failedRun.ID)
	if err != nil {
		return "", err
	}

	for _, task := range run.Tasks {
		if task.Status == agola.RunTaskStatusFailed {
			taskFailed, err := agola.GetTask(run.ID, task.ID)
			if err != nil {
				return "", err
			}

			if taskFailed.SetupStep.Phase == agola.ExecutorTaskPhaseFailed {
				logs, err := agola.GetLogs(run.ID, task.ID, -1)
				if err != nil {
					return "", err
				}

				body += "\n\n#Task setup " + task.Name + " failed\n" + logs
			}

			for stepID, step := range taskFailed.Steps {
				if step.Phase == agola.ExecutorTaskPhaseFailed {

					logs, err := agola.GetLogs(run.ID, task.ID, stepID)
					if err != nil {
						return "", err
					}

					body += "\n\n#task " + task.Name + " #step " + step.Name + "\n" + logs
				}
			}
		}
	}

	log.Println("* mail body *", body)

	return body, nil
}

func CheckIfNewRunsPresent(project *model.Project) bool {
	lastRun := project.GetLastRun()
	runList, _ := agola.GetRuns(project.AgolaProjectID, true, "finished", nil, 1, true)

	return runList != nil && len(*runList) != 0 && (*runList)[0].StartTime.After(lastRun.RunStartDate)
}

func findGiteaUsersEmailRepositoryOwner(gitSource *model.GitSource, organizationName string, gitRepoPath string) (*[]string, error) {
	retVal := make([]string, 0)

	teams, err := gitea.GetRepositoryTeams(gitSource, organizationName, gitRepoPath)
	if err != nil {
		return nil, err
	}

	for _, team := range *teams {
		if strings.Compare(team.Permission, "owner") != 0 {
			continue
		}

		users, err := gitea.GetTeamMembers(gitSource, team.ID)
		if err != nil {
			continue
		}

		for _, user := range *users {
			retVal = append(retVal, user.Email)
		}
	}

	return &retVal, nil
}

func findGithubUsersRepositoryOwner(gitSource *model.GitSource, organizationName string, gitRepoPath string) (*[]string, error) {
	retVal := make([]string, 0)

	users, err := github.GetRepositoryMembers(gitSource, organizationName, gitRepoPath)
	if err != nil {
		return nil, err
	}

	for _, user := range *users {
		if strings.Compare(user.Role, "owner") == 0 {
			retVal = append(retVal, user.Email)
		}
	}

	return &retVal, nil
}

func getEmailByRun(run *agola.RunDto, gitSource *model.GitSource, organizationName string, gitRepoPath string) []string {
	retVal := make([]string, 0)

	commitMetadata, err := git.GetCommitMetadata(gitSource, organizationName, gitRepoPath, run.GetCommitSha())
	if err == nil && commitMetadata != nil {
		retVal = append(retVal, commitMetadata.GetAuthorEmail())

		if commitMetadata.Parents != nil {
			for _, parent := range commitMetadata.Parents {
				commitParentMetadata, err := git.GetCommitMetadata(gitSource, organizationName, gitRepoPath, parent.Sha)
				if err == nil && commitParentMetadata != nil {
					retVal = append(retVal, commitParentMetadata.GetAuthorEmail())
				}
			}
		}
	}

	return retVal
}

//Take only the run by webhook, discard others(for example directrun)
func takeWebhookTrigger(runs *[]agola.RunDto) *[]agola.RunDto {
	retVal := make([]agola.RunDto, 0)

	if runs != nil {
		for _, run := range *runs {
			if run.IsWebhookCreationTrigger() {
				retVal = append(retVal, run)
			}
		}
	}

	return &retVal
}

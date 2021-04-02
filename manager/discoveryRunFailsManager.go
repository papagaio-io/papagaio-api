package manager

import (
	"fmt"
	"log"
	"strings"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

func discoveryRunFails(db repository.Database) {
	for {
		log.Println("Start discoveryRunFails")

		organizations, _ := db.GetOrganizations()
		for _, org := range *organizations {
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
					log.Println("no new runs found...")
					continue
				}

				//If there are new runs asks for other runs
				runList, _ := agola.GetRuns(project.AgolaProjectID, false, "finished", &project.OlderRunFaild.ID, 0, true)
				project.LastRun = model.RunInfo{ID: (*runList)[len(*runList)-1].ID, RunStartDate: (*runList)[len(*runList)-1].StartTime, Branch: (*runList)[len(*runList)-1].GetBranchName()}

				//I use the save the faild runs by branch, next I take the older and save to db for optimize agola.GetRuns
				runsFaildMap := make(map[string]agola.RunDto)

				runList = takeWebhookTrigger(runList)

				branchRunLists := subdivideRunsByBranch(runList)
				for branch, branchRunList := range *branchRunLists {
					//Prendo l'ultima run del branch dal db, se c'è
					//Se c'è tolgo le run precedenti di agola da a questa ultima
					lastFailesRunSaved, ok := project.LastBranchRunFailsMap[branch]
					if ok {
						branchRunList = deleteOlderRunsBy(&branchRunList, lastFailesRunSaved)
					}

					for _, run := range branchRunList {
						if run.Result == agola.RunResultFailed && run.StartTime.After(lastFailesRunSaved.RunStartDate) { //Se la prima run fallita di agola corrisponde a quella presa dal db suppongo di avere già notificato gli utenti al polling precedente
							log.Println("Found run failed!")

							project.LastBranchRunFailsMap[branch] = model.RunInfo{ID: run.ID, RunStartDate: run.StartTime, Branch: branch}
							runsFaildMap[branch] = run

							emailMap := getUsersEmailMap(gitSource, &org, project.GitRepoPath, run)

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

					olderFailedRun := findOlderRun(&runsFaildMap)
					if olderFailedRun != nil {
						project.OlderRunFaild = model.RunInfo{ID: olderFailedRun.ID, RunStartDate: olderFailedRun.StartTime, Branch: branch}
					}
				}
				org.Projects[projectName] = project
			}
			db.SaveOrganization(&org)
		}

		log.Println("End discoveryRunFails")
		time.Sleep(5 * time.Minute)
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

const bodyMainTemplate string = `[%s/%s] FIX Agola Run (#%s)\nSee: <a href="%s"` + ">click here</a>"
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
	body := fmt.Sprintf(bodyMainTemplate, organizationName, projectName, fmt.Sprint(failedRun.Counter), runUrl)

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

func findOlderRun(runs *map[string]agola.RunDto) *agola.RunDto {
	if runs == nil || len(*runs) == 0 {
		return nil
	}

	var retVal *agola.RunDto = nil
	for _, run := range *runs {
		if retVal == nil || run.StartTime.Before(retVal.StartTime) {
			retVal = &run
		}
	}

	return retVal
}

func CheckIfNewRunsPresent(project *model.Project) bool {
	if project.LastBranchRunFailsMap == nil {
		project.LastBranchRunFailsMap = make(map[string]model.RunInfo)
	}

	isLastRun := !project.LastRun.RunStartDate.IsZero()
	runList, _ := agola.GetRuns(project.AgolaProjectID, isLastRun, "finished", nil, 1, true)

	return runList != nil && len(*runList) != 0 && (*runList)[0].StartTime.After(project.LastRun.RunStartDate)
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

func deleteOlderRunsBy(runs *[]agola.RunDto, firstRun model.RunInfo) []agola.RunDto {
	if runs == nil {
		return nil
	}

	retVal := make([]agola.RunDto, 0)
	for _, run := range *runs {
		if run.StartTime.Equal(firstRun.RunStartDate) || run.StartTime.After(firstRun.RunStartDate) {
			retVal = append(retVal, run)
		}
	}

	return retVal
}

func subdivideRunsByBranch(runs *[]agola.RunDto) *map[string][]agola.RunDto {
	retVal := make(map[string][]agola.RunDto)

	for _, run := range *runs {
		branch := run.GetBranchName()
		if _, ok := retVal[branch]; !ok {
			retVal[branch] = make([]agola.RunDto, 1)
			retVal[branch][0] = run
		} else {
			retVal[branch] = append(retVal[branch], run)
		}
	}

	return &retVal
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

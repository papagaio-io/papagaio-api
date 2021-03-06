package trigger

import (
	"fmt"
	"log"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/trigger/dto"
	"wecode.sorint.it/opensource/papagaio-api/types"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

func StartRunFailsDiscovery(db repository.Database, tr utils.ConfigUtils, commonMutex *utils.CommonMutex, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway, rtDto *dto.TriggerRunTimeDto) {
	go discoveryRunFails(db, tr, commonMutex, agolaApi, gitGateway, rtDto)
}

/*
Scan Agola project runs and store it for elaborating of reports.
If find failed runs send email to users
*/
func discoveryRunFails(db repository.Database, tr utils.ConfigUtils, commonMutex *utils.CommonMutex, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway, rtDto *dto.TriggerRunTimeDto) {
	for {
		rtDto.IsRunning = true
		rtDto.TimerLastRun = time.Now()

		log.Println("Start discoveryRunFails")
		rtDto.TriggerLastRun = time.Now()

		organizationsRef, _ := db.GetOrganizationsRef()

		for _, organizationRef := range organizationsRef {
			mutex := utils.ReserveOrganizationMutex(organizationRef, commonMutex)
			mutex.Lock()

			org, _ := db.GetOrganizationByAgolaRef(organizationRef)
			if org == nil {
				log.Println("discoveryRunFails organization ", organizationRef, "not found")

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

				continue
			}

			gitSource, err := db.GetGitSourceByName(org.GitSourceName)
			if gitSource == nil || err != nil || org.Projects == nil {
				log.Println("discoveryRunFails gitsource not fount for", organizationRef, "organization")

				mutex.Unlock()
				utils.ReleaseOrganizationMutex(organizationRef, commonMutex)

				continue
			}

			user, _ := db.GetUserByUserId(org.UserIDConnected)

			for projectName, project := range org.Projects {
				if project.Archivied {
					continue
				}

				checkNewRuns := CheckIfNewRunsPresent(&project, agolaApi)
				if !checkNewRuns {
					log.Println("no new runs found for project", projectName)
					continue
				}

				//If there are new runs asks for other runs
				lastRun := project.GetLastRun()
				runList, _ := agolaApi.GetRuns(project.AgolaProjectID, false, "finished", &lastRun.ID, 0, true)

				runList = takeWebhookTrigger(runList)

				for _, run := range *runList {
					if !run.IsBranch() { //skip tags
						continue
					}

					newRun := model.RunInfo{
						ID:           run.ID,
						Branch:       run.GetBranchName(),
						RunStartDate: run.StartTime,
						RunEndDate:   run.EndTime,
						Phase:        types.RunPhase(run.Phase),
						Result:       types.RunResult(run.Result),
					}
					project.PushNewRun(newRun)

					//

					if run.Result == agola.RunResultFailed && run.StartTime.After(lastRun.RunStartDate) {
						log.Println("Found run failed!")
						emailMap := getUsersEmailMap(gitSource, user, org, project.GitRepoPath, run, gitGateway)
						log.Println("send emails to:", emailMap)

						body, err := makeBody(org, project.GitRepoPath, run, agolaApi)
						if err != nil {
							log.Println("Failed to make email body")
							continue
						}
						subject := makeSubject(org, project.GitRepoPath, run)

						if utils.CanSendEmail() {
							utils.SendConfirmEmail(emailMap, nil, subject, body)
						} else {
							log.Println("Can not send email, settings are not correct")
						}
					}
				}

				org.Projects[projectName] = project
			}
			err = db.SaveOrganization(org)

			if err != nil {
				log.Println("error in SaveOrganization:", err)
			}

			mutex.Unlock()
			utils.ReleaseOrganizationMutex(organizationRef, commonMutex)
		}

		rtDto.IsRunning = false

		fmt.Println("discoveryRunFails end")

		select {
		case message := <-rtDto.Chan:
			fmt.Println("discoveryRunFails message:", message)

		case <-time.After(time.Duration(time.Minute.Nanoseconds() * int64(tr.GetRunFailedTriggerTime()))):
		}
	}
}

func getUsersEmailMap(gitSource *model.GitSource, user *model.User, organization *model.Organization, gitRepoPath string, failedRun agola.RunDto, gitGateway *git.GitGateway) map[string]bool {
	emails := make(map[string]bool)

	//Find all users that commited the failed run and parents
	emailUsersCommitted := getEmailByRun(&failedRun, gitSource, user, organization.GitPath, gitRepoPath, gitGateway)

	//Users owner of the organization and users owner of the repository
	var usersRepoOwners *[]string

	usersRepoOwners, _ = gitGateway.GetEmailsRepositoryUsersOwner(gitSource, user, organization.GitPath, gitRepoPath)

	for _, email := range emailUsersCommitted {
		emails[email] = true
	}

	if usersRepoOwners != nil {
		for _, email := range *usersRepoOwners {
			emails[email] = true
		}
	}

	if organization.ExternalUsers != nil {
		for email := range organization.ExternalUsers {
			emails[email] = true
		}
	}

	return emails
}

const bodyMessageTemplate string = "[%s/%s] FIX Agola Run (#%s)\n"
const bodyLinkTemplate string = `See: <a href="%s">click here</a>`
const subjectTemplate string = "Run failed in Agola: %s ?? %s ?? release #%s"
const runAgolaPath string = "%s/org/%s/projects/%s.proj/runs/%s"

func makeSubject(organization *model.Organization, projectName string, failedRun agola.RunDto) string {
	return fmt.Sprintf(subjectTemplate, organization.GitPath, projectName, fmt.Sprint(failedRun.Counter))
}

func getRunAgolaUrl(organization *model.Organization, projectName string, runID string) string {
	return fmt.Sprintf(runAgolaPath, config.Config.Agola.AgolaAddr, organization.AgolaOrganizationRef, projectName, runID)
}

func makeBody(organization *model.Organization, projectName string, failedRun agola.RunDto, agolaApi agola.AgolaApiInterface) (string, error) {
	runUrl := getRunAgolaUrl(organization, projectName, failedRun.ID)
	body := fmt.Sprintf(bodyMessageTemplate, organization.GitPath, projectName, fmt.Sprint(failedRun.Counter))
	body += fmt.Sprintf(bodyLinkTemplate, runUrl)

	run, err := agolaApi.GetRun(failedRun.ID)
	if err != nil {
		return "", err
	}

	for _, task := range run.Tasks {
		if task.Status == agola.RunTaskStatusFailed {
			taskFailed, err := agolaApi.GetTask(run.ID, task.ID)
			if err != nil {
				return "", err
			}

			if taskFailed.SetupStep.Phase == agola.ExecutorTaskPhaseFailed {
				logs, err := agolaApi.GetLogs(run.ID, task.ID, -1)
				if err != nil {
					return "", err
				}

				body += "\n\n#Task setup " + task.Name + " failed\n" + logs
			}

			for stepID, step := range taskFailed.Steps {
				if step.Phase == agola.ExecutorTaskPhaseFailed {

					logs, err := agolaApi.GetLogs(run.ID, task.ID, stepID)
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

func CheckIfNewRunsPresent(project *model.Project, agolaApi agola.AgolaApiInterface) bool {
	lastRun := project.GetLastRun()
	runList, _ := agolaApi.GetRuns(project.AgolaProjectID, true, "finished", nil, 1, false)

	return runList != nil && len(*runList) != 0 && (*runList)[0].StartTime.After(lastRun.RunStartDate)
}

func getEmailByRun(run *agola.RunDto, gitSource *model.GitSource, user *model.User, organizationName string, gitRepoPath string, gitGateway *git.GitGateway) []string {
	retVal := make([]string, 0)

	commitMetadata, err := gitGateway.GetCommitMetadata(gitSource, user, organizationName, gitRepoPath, run.GetCommitSha())
	if err == nil && commitMetadata != nil {
		email := commitMetadata.GetAuthorEmail()
		if email != nil {
			retVal = append(retVal, *email)
		}

		if commitMetadata.Parents != nil {
			for _, parent := range commitMetadata.Parents {
				commitParentMetadata, err := gitGateway.GetCommitMetadata(gitSource, user, organizationName, gitRepoPath, parent.Sha)
				if err == nil && commitParentMetadata != nil {
					email = commitParentMetadata.GetAuthorEmail()
					if email != nil {
						retVal = append(retVal, *email)
					}
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

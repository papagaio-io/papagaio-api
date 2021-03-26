package manager

import (
	"log"
	"strings"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
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
				continue
			}

			for projectName, project := range org.Projects {
				if project.Archivied {
					continue
				}

				if project.LastBranchRunFailsMap == nil {
					project.LastBranchRunFailsMap = make(map[string]model.RunInfo)
				}

				olderDbBranch := getOlderBranchRun(&project.LastBranchRunFailsMap)
				isLastRun := olderDbBranch != nil

				var startRunID *string = nil
				var olderDbRunInfo model.RunInfo
				if olderDbBranch != nil {
					olderDbRunInfo = project.LastBranchRunFailsMap[*olderDbBranch]
					startRunID = &olderDbRunInfo.RunID
				}

				runList, err := agola.GetRuns(project.AgolaProjectID, isLastRun, "finished", startRunID, 1, true)

				if err != nil || runList == nil || len(*runList) == 0 {
					log.Println("no runs found...")
					continue
				}

				//Skip if there are no new runs
				if startRunID != nil { //so olderDbRunInfo is not empty
					if !(*runList)[len(*runList)-1].StartTime.After(project.LastRun.RunStartDate) {
						log.Println("no new runs found...")
						continue
					}
				}

				//If there are new runs asks for other runs
				runList, err = agola.GetRuns(project.AgolaProjectID, false, "finished", startRunID, 0, true)
				project.LastRun = model.RunInfo{RunID: (*runList)[len(*runList)-1].ID, RunStartDate: *(*runList)[len(*runList)-1].StartTime}

				runList = takeWebhookTrigger(runList)

				runListSubdivided := subdivideRunsByBranch(runList)
				for branch, branchRunList := range *runListSubdivided {
					//Prendo l'ultima run del branch dal db, se c'è
					//Se c'è tolgo le run precedenti di agola da a questa ultima
					lastFailesRunSaved, ok := project.LastBranchRunFailsMap[branch]
					if ok {
						branchRunList = deleteOlderRunsBy(&branchRunList, lastFailesRunSaved)
					}

					runToSaveOnDb := lastFailesRunSaved
					successRun := make([]agolaApi.RunDto, 0)

					for _, run := range branchRunList {
						if run.Result == agola.RunResultSuccess {
							successRun = append(successRun, run)
						} else if run.Result == agolaApi.RunResultFailed && run.StartTime.After(lastFailesRunSaved.RunStartDate) { //Se la prima run fallita di agola corrisponde a quella presa dal db suppongo di avere già notificato gli utenti al polling precedente
							log.Println("Found run failed!")

							runToSaveOnDb = model.RunInfo{RunID: run.ID, RunStartDate: *run.StartTime}

							//Find all users that commited the failed run and success runs
							localRuns := append(successRun, run)
							emailUsersCommitted := getEmailByRuns(&localRuns, gitSource, org.Name, project.GitRepoPath)

							//After email is sent we empty successRun
							successRun = make([]agolaApi.RunDto, 0)

							//Users owner of the organization and users owner of the repository
							var usersRepoOwners *[]string
							if org.OtherUserToNotify != nil {
								for _, email := range org.OtherUserToNotify {
									emailUsersCommitted = append(emailUsersCommitted, email)
								}
							}

							if gitSource.GitType == model.Gitea {
								usersRepoOwners, _ = findGiteaUsersEmailRepositoryOwner(gitSource, org.Name, project.GitRepoPath)
							} else {
								usersRepoOwners, _ = findGithubUsersRepositoryOwner(gitSource, org.Name, project.GitRepoPath)
							}

							//Create map without duplicate email
							emailMap := make(map[string]bool)

							for _, email := range emailUsersCommitted {
								emailMap[email] = true
							}

							if usersRepoOwners != nil {
								for _, email := range *usersRepoOwners {
									emailMap[email] = true
								}
							}

							if org.OtherUserToNotify != nil {
								for _, email := range org.OtherUserToNotify {
									emailMap[email] = true
								}
							}

							log.Println("send emails to:", emailMap)

							sendConfirmEmail(emailMap, nil, "Run Agola faild subject", "Run Agola faild body")
						}
					}

					project.LastBranchRunFailsMap[branch] = runToSaveOnDb
				}
				org.Projects[projectName] = project
			}
			db.SaveOrganization(&org)
		}

		log.Println("End discoveryRunFails")
		time.Sleep(3 * time.Minute)
	}
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

func getEmailByRuns(runs *[]agolaApi.RunDto, gitSource *model.GitSource, organizationName string, gitRepoPath string) []string {
	retVal := make([]string, 0)

	for _, run := range *runs {
		commitMetadata, err := git.GetCommitMetadata(gitSource, organizationName, gitRepoPath, run.GetCommitSha())
		if commitMetadata == nil || err != nil {
			continue
		}
		retVal = append(retVal, commitMetadata.GetAuthorEmail())
	}

	return retVal
}

func deleteOlderRunsBy(runs *[]agolaApi.RunDto, firstRun model.RunInfo) []agolaApi.RunDto {
	if runs == nil {
		return nil
	}

	retVal := make([]agolaApi.RunDto, 0)
	for _, run := range *runs {
		if run.StartTime.Equal(firstRun.RunStartDate) || run.StartTime.After(firstRun.RunStartDate) {
			retVal = append(retVal, run)
		}
	}

	return retVal
}

func subdivideRunsByBranch(runs *[]agolaApi.RunDto) *map[string][]agolaApi.RunDto {
	retVal := make(map[string][]agolaApi.RunDto)

	for _, run := range *runs {
		branch := run.GetBranchName()
		if _, ok := retVal[branch]; !ok {
			retVal[branch] = make([]agolaApi.RunDto, 1)
			retVal[branch][0] = run
		} else {
			retVal[branch] = append(retVal[branch], run)
		}
	}

	return &retVal
}

func getOlderBranchRun(runList *map[string]model.RunInfo) *string {
	if runList == nil || len(*runList) == 0 {
		return nil
	}

	olderRunTime := time.Now()
	var retVal *string
	for branch, run := range *runList {
		if run.RunStartDate.Before(olderRunTime) {
			olderRunTime = run.RunStartDate
			retVal = &branch
		}
	}

	return retVal
}

//Take only the run by webhook, discard others(for example directrun)
func takeWebhookTrigger(runs *[]agolaApi.RunDto) *[]agolaApi.RunDto {
	retVal := make([]agolaApi.RunDto, 0)

	if runs != nil {
		for _, run := range *runs {
			if run.IsWebhookCreationTrigger() {
				retVal = append(retVal, run)
			}
		}
	}

	return &retVal
}

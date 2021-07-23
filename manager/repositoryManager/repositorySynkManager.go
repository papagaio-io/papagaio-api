package repositoryManager

import (
	"log"
	"strings"

	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

//Inserisco tutti i repository di git su agola
func CheckoutAllGitRepository(db repository.Database, user *model.User, organization *model.Organization, gitSource *model.GitSource, agolaApi agolaApi.AgolaApiInterface, gitGateway *git.GitGateway) {
	log.Println("Start AddAllGitRepository")

	repositoryList, _ := gitGateway.GetRepositories(gitSource, user, organization.GitPath)

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	for _, repo := range *repositoryList {
		if !utils.EvaluateBehaviour(organization, repo) {
			continue
		}

		log.Println("Start add repository:", repo)

		agolaConfExists, _ := gitGateway.CheckRepositoryAgolaConfExists(gitSource, user, organization.GitPath, repo)
		project := model.Project{GitRepoPath: repo, Archivied: true, AgolaProjectRef: utils.ConvertToAgolaProjectRef(repo)}

		if agolaConfExists {
			projectID, err := agolaApi.CreateProject(repo, project.AgolaProjectRef, organization, gitSource.AgolaRemoteSource, user)
			project.AgolaProjectID = projectID
			project.Archivied = false
			if err != nil {
				log.Println("Warning!!! Agola CreateProject API error:", err.Error())
			}
		}

		organization.Projects[repo] = project
		err := db.SaveOrganization(organization)
		if err != nil {
			log.Println("error in SaveOrganization:", err)
			return
		}

		BranchSynck(db, user, gitSource, organization, repo, gitGateway)

		log.Println("End add repository:", repo)
	}

	log.Println("End CheckoutAllGitRepository")
}

func SynkGitRepositorys(db repository.Database, user *model.User, organization *model.Organization, gitSource *model.GitSource, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway) error {
	log.Println("Start SynkGitRepositorys for", organization.GitPath)

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	gitRepositoryList, err := gitGateway.GetRepositories(gitSource, user, organization.GitPath)
	if err == nil {
		log.Println("git GetRepositories err:", err)
	}

	if gitRepositoryList != nil {
		for projectName, project := range organization.Projects {
			log.Println("SynkGitRepositorys git repository:", projectName)
			gitRepoExists := false
			for _, gitRepo := range *gitRepositoryList {
				if strings.Compare(projectName, gitRepo) == 0 {
					gitRepoExists = true
					break
				}
			}
			if !gitRepoExists {
				err := agolaApi.DeleteProject(organization, project.AgolaProjectRef, user)
				if err == nil {
					delete(organization.Projects, projectName)
				} else {
					log.Println("Agola DeleteProject error:", err)
				}
			} else {
				agolaExists, agolaProjectID := agolaApi.CheckProjectExists(organization, project.AgolaProjectRef)
				if !agolaExists && !project.Archivied {
					delete(organization.Projects, projectName)
				} else {
					project.AgolaProjectID = agolaProjectID
					organization.Projects[projectName] = project
				}
			}
		}

		for _, repo := range *gitRepositoryList {
			if !utils.EvaluateBehaviour(organization, repo) {
				delete(organization.Projects, repo)

				agolaProjectRef := utils.ConvertToAgolaProjectRef(repo)
				if exists, _ := agolaApi.CheckProjectExists(organization, agolaProjectRef); exists {
					err := agolaApi.DeleteProject(organization, agolaProjectRef, user)
					if err != nil {
						log.Println("Agola DeleteProject error:", err)
					}
				}

				continue
			}

			var project model.Project
			if p, ok := organization.Projects[repo]; !ok {
				project = model.Project{GitRepoPath: repo, AgolaProjectRef: utils.ConvertToAgolaProjectRef(repo)}
				organization.Projects[repo] = project
			} else {
				project = p
			}

			BranchSynck(db, user, gitSource, organization, repo, gitGateway)

			agolaConfExists, _ := gitGateway.CheckRepositoryAgolaConfExists(gitSource, user, organization.GitPath, repo)
			if !agolaConfExists {
				if project, ok := organization.Projects[repo]; ok && !project.Archivied {
					err := agolaApi.ArchiveProject(organization, project.AgolaProjectRef)
					if err == nil {
						project.Archivied = true
						organization.Projects[repo] = project
					}
				}

				continue
			}

			if exists, projectID := agolaApi.CheckProjectExists(organization, utils.ConvertToAgolaProjectRef(repo)); exists {
				if project, ok := organization.Projects[repo]; ok {
					project.AgolaProjectID = projectID
					if project.Archivied {
						err := agolaApi.UnarchiveProject(organization, utils.ConvertToAgolaProjectRef(repo))
						if err == nil {
							project.Archivied = false
							organization.Projects[repo] = project
						}
					}
					organization.Projects[repo] = project
				}

				continue
			}

			log.Println("Start add repository:", repo)
			projectID, err := agolaApi.CreateProject(repo, utils.ConvertToAgolaProjectRef(repo), organization, gitSource.AgolaRemoteSource, user)
			if err != nil {
				log.Println("Warning!!! Agola CreateProject API error:", err.Error())
				break
			}
			project.AgolaProjectID = projectID
			organization.Projects[repo] = project
			log.Println("End add repository:", repo)
		}
	}

	err = db.SaveOrganization(organization)
	if err != nil {
		return err
	}

	log.Println("End SynkGitRepositorys for", organization.GitPath)

	return nil
}

func BranchSynck(db repository.Database, user *model.User, gitSource *model.GitSource, organization *model.Organization, repositoryName string, gitGateway *git.GitGateway) {
	if _, exists := organization.Projects[repositoryName]; !exists {
		return
	}

	if organization.Projects[repositoryName].Branchs == nil {
		project := organization.Projects[repositoryName]
		project.Branchs = make(map[string]model.Branch)
		organization.Projects[repositoryName] = project
	}

	branchList := gitGateway.GetBranches(gitSource, user, organization.GitPath, repositoryName)

	for branch := range branchList {
		if _, ok := organization.Projects[repositoryName].Branchs[branch]; !ok {
			organization.Projects[repositoryName].Branchs[branch] = model.Branch{Name: branch}
		}
	}

	for branch := range organization.Projects[repositoryName].Branchs {
		if _, ok := branchList[branch]; !ok {
			delete(organization.Projects[repositoryName].Branchs, branch)
		}
	}

	err := db.SaveOrganization(organization)
	if err != nil {
		log.Println("error in SaveOrganization:", err)
	}
}

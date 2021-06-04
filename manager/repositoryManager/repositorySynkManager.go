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

	repositoryList, _ := gitGateway.GetRepositories(gitSource, user, organization.Name)

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	for _, repo := range *repositoryList {
		if !utils.EvaluateBehaviour(organization, repo) {
			continue
		}

		log.Println("Start add repository:", repo)

		agolaConfExists, _ := gitGateway.CheckRepositoryAgolaConfExists(gitSource, user, organization.Name, repo)
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
		db.SaveOrganization(organization)

		BranchSynck(db, user, gitSource, organization, repo, gitGateway)

		log.Println("End add repository:", repo)
	}

	log.Println("End CheckoutAllGitRepository")
}

func SynkGitRepositorys(db repository.Database, user *model.User, organization *model.Organization, gitSource *model.GitSource, agolaApi agola.AgolaApiInterface, gitGateway *git.GitGateway) error {
	log.Println("Start SynkGitRepositorys for", organization.Name)

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	gitRepositoryList, _ := gitGateway.GetRepositories(gitSource, user, organization.Name)

	for projectName, project := range organization.Projects {
		gitRepoExists := false
		for _, gitRepo := range *gitRepositoryList {
			if strings.Compare(projectName, gitRepo) == 0 {
				gitRepoExists = true
				break
			}
		}
		if !gitRepoExists {
			agolaApi.DeleteProject(organization, project.AgolaProjectRef, user)
			delete(organization.Projects, projectName)
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
				agolaApi.DeleteProject(organization, agolaProjectRef, user)
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

		agolaConfExists, _ := gitGateway.CheckRepositoryAgolaConfExists(gitSource, user, organization.Name, repo)
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

	db.SaveOrganization(organization)

	log.Println("End SynkGitRepositorys for", organization.Name)

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

	branchList := gitGateway.GetBranches(gitSource, user, organization.Name, repositoryName)

	for branch, _ := range branchList {
		if _, ok := organization.Projects[repositoryName].Branchs[branch]; !ok {
			organization.Projects[repositoryName].Branchs[branch] = model.Branch{Name: branch}
		}
	}

	for branch, _ := range organization.Projects[repositoryName].Branchs {
		if _, ok := branchList[branch]; !ok {
			delete(organization.Projects[repositoryName].Branchs, branch)
		}
	}

	db.SaveOrganization(organization)
}

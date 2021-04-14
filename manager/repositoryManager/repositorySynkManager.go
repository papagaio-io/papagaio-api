package repositoryManager

import (
	"errors"
	"log"
	"strings"

	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	gitApi "wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

//Inserisco tutti i repository di git su agola
func CheckoutAllGitRepository(db repository.Database, organization *model.Organization, gitSource *model.GitSource) error {
	log.Println("Start AddAllGitRepository")

	repositoryList, _ := gitApi.GetRepositories(gitSource, organization.Name)
	log.Println("repositoryList:", *repositoryList)

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	for _, repo := range *repositoryList {
		if !utils.EvaluateBehaviour(organization, repo) {
			continue
		}

		agolaConfExists, _ := git.CheckRepositoryAgolaConf(gitSource, organization.Name, repo)
		if !agolaConfExists {
			continue
		}

		log.Println("Start add repository:", repo)

		projectID, err := agolaApi.CreateProject(repo, organization, gitSource.AgolaRemoteSource, gitSource.AgolaToken)

		if err != nil {
			log.Println("Warning!!! Agola CreateProject API error:", err.Error())
			return errors.New(err.Error())
		}

		project := model.Project{OrganizationID: organization.ID, GitRepoPath: repo, AgolaProjectID: projectID}
		organization.Projects[repo] = project
		db.SaveOrganization(organization)

		log.Println("End add repository:", repo)
	}

	log.Println("End AddAllGitRepository")

	return nil
}

func SynkGitRepositorys(db repository.Database, organization *model.Organization, gitSource *model.GitSource) error {
	log.Println("Start SynkGitRepositorys for", organization.Name)

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	gitRepositoryList, _ := gitApi.GetRepositories(gitSource, organization.Name)

	//if some project is not present in agola I remove from db
	for projectName, project := range organization.Projects {
		agolaExists, agolaProjectID := agolaApi.CheckProjectExists(organization.Name, projectName)

		if !agolaExists {
			delete(organization.Projects, projectName)
			continue
		}

		gitRepoExists := false
		for _, gitRepo := range *gitRepositoryList {
			if strings.Compare(project.GitRepoPath, gitRepo) == 0 {
				gitRepoExists = true
				break
			}
		}
		if !gitRepoExists {
			agolaApi.DeleteProject(organization.Name, project.GitRepoPath, gitSource.AgolaToken)
			delete(organization.Projects, project.GitRepoPath)

		} else {
			project.AgolaProjectID = agolaProjectID
			organization.Projects[projectName] = project
		}
	}

	for _, repo := range *gitRepositoryList {
		if !utils.EvaluateBehaviour(organization, repo) {
			if _, ok := organization.Projects[repo]; ok {
				delete(organization.Projects, repo)
			}

			if exists, _ := agolaApi.CheckProjectExists(organization.Name, repo); exists {
				agolaApi.DeleteProject(organization.Name, repo, gitSource.AgolaToken)
			}

			continue
		}

		BranchSynck(db, gitSource, organization, repo)

		agolaConfExists, _ := git.CheckRepositoryAgolaConf(gitSource, organization.Name, repo)
		if !agolaConfExists {
			if project, ok := organization.Projects[repo]; ok && !project.Archivied {
				err := agolaApi.ArchiveProject(organization.Name, repo)
				if err != nil {
					project.Archivied = true
					organization.Projects[repo] = project
				}
			}

			continue
		}

		if exists, projectID := agolaApi.CheckProjectExists(organization.Name, repo); exists {
			if project, ok := organization.Projects[repo]; ok {
				project.AgolaProjectID = projectID
				if project.Archivied {
					err := agolaApi.UnarchiveProject(organization.Name, repo)
					if err != nil {
						project.Archivied = false
						organization.Projects[repo] = project
					}
				}
				organization.Projects[repo] = project
			}

			continue
		}

		log.Println("Start add repository:", repo)

		projectID, err := agolaApi.CreateProject(repo, organization, gitSource.AgolaRemoteSource, gitSource.AgolaToken)

		if err != nil {
			log.Println("Warning!!! Agola CreateProject API error:", err.Error())
			break
		}

		project := model.Project{OrganizationID: organization.ID, GitRepoPath: repo, AgolaProjectID: projectID}
		organization.Projects[repo] = project

		log.Println("End add repository:", repo)
	}

	db.SaveOrganization(organization)

	log.Println("End SynkGitRepositorys for", organization.Name)

	return nil
}

func BranchSynck(db repository.Database, gitSource *model.GitSource, organization *model.Organization, repositoryName string) {
	if organization.Projects[repositoryName].Branchs == nil {
		project := organization.Projects[repositoryName]
		project.Branchs = make(map[string]model.Branch)
		organization.Projects[repositoryName] = project
	}

	branchList := git.GetBranches(gitSource, organization.Name, repositoryName)

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

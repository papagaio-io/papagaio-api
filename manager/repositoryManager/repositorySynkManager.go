package repositoryManager

import (
	"errors"
	"log"

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
	log.Println("Start SynkGitRepositorys")

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	gitRepositoryList, _ := gitApi.GetRepositories(gitSource, organization.Name)
	log.Println("repositoryList:", *gitRepositoryList)

	//if some project is not present in agola I remove from db
	for _, project := range organization.Projects {
		if !agolaApi.CheckProjectExists(organization.Name, project.GitRepoPath) {
			delete(organization.Projects, project.GitRepoPath)
		}
	}

	for _, repo := range *gitRepositoryList {
		if !utils.EvaluateBehaviour(organization, repo) {
			if _, ok := organization.Projects[repo]; ok {
				delete(organization.Projects, repo)
			}

			if agolaApi.CheckProjectExists(organization.Name, repo) {
				agolaApi.DeleteProject(organization.Name, repo, gitSource.AgolaToken)
			}

			continue
		}

		agolaConfExists, _ := git.CheckRepositoryAgolaConf(gitSource, organization.Name, repo)
		if !agolaConfExists {
			if project, ok := organization.Projects[repo]; ok && !project.Archivied {
				agolaApi.ArchiveProject(organization.Name, repo)
				project.Archivied = true
			}

			continue
		}

		if agolaApi.CheckProjectExists(organization.Name, repo) {
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

		log.Println("End add repository:", repo)
	}

	db.SaveOrganization(organization)

	log.Println("End SynkGitRepositorys")

	return nil
}

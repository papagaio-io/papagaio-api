package repositoryManager

import (
	"fmt"

	agolaApi "wecode.sorint.it/opensource/papagaio-api/api/agola"
	gitApi "wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/utils"
)

//Inserisco tutti i repository di git su agola
func AddAllGitRepository(db repository.Database, organization *model.Organization, gitSource *model.GitSource) {
	repository, _ := gitApi.GetRepositories(gitSource, organization.Name)

	if organization.Projects == nil {
		organization.Projects = make(map[string]model.Project)
	}

	for _, repo := range *repository {
		if !utils.EvaluateBehaviour(organization, repo) {
			continue
		}

		gitSource, _ := db.GetGitSourceById(organization.GitSourceID)
		projectID, err := agolaApi.CreateProject(repo, organization, gitSource.AgolaRemoteSource, gitSource.AgolaToken)

		if err != nil {
			fmt.Println("Warning!!! Agola CreateProject API error!")
			return
		}

		project := model.Project{OrganizationID: organization.ID, GitRepoPath: repo, AgolaProjectID: projectID}
		organization.Projects[repo] = project
		db.SaveOrganization(organization)
	}
}

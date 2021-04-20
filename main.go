package main

import (
	"fmt"
	"time"

	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
	"wecode.sorint.it/opensource/papagaio-api/cmd"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/manager/repositoryManager"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

func main() {
	//testRunFails()
	//testGitHub()

	//testSynkWecode()

	cmd.Execute()

	/*db := repository.NewAppDb(config.Config)

	gs, _ := db.GetGitSourceByName("wecodedev")
	fmt.Println("gs:", gs)*/
}

func testRunFails() {
	/*db := repository.NewAppDb(config.Config)

	gitSource, _ := db.GetGitSourceByName("test")
	if gitSource == nil {
		gitSource = &model.GitSource{Name: "test", GitType: model.Gitea, GitToken: "7171eee0d8404600d4da7836f81c07df5a43daaf", GitAPIURL: "https://wecode.sorintdev.it", AgolaRemoteSource: "alessandropinnawecode", AgolaToken: "5adacf4e054f04133c60ac7a4f139dc9aa85b232"}
		db.SaveGitSource(gitSource)
	}

	projects := make(map[string]model.Project)
	projects["ProvaRepoForked"] = model.Project{GitRepoPath: "ProvaRepoForked", AgolaProjectID: "e90211dd-6801-4ce5-a59f-58dcc6eb4391", Archivied: false}

	organization := model.Organization{Name: "Sorint", GitSourceID: gitSource.ID, Projects: projects}
	db.SaveOrganization(&organization)

	trigger.StartRunFailsDiscovery(&db)*/
}

func testGitHub() {
	githubSource := &model.GitSource{GitToken: "1af89575f328da999fffdd040667dc7b035e05ee"}

	/*list, _ := github.GetRepositories(githubSource, "Sorinttest")
	fmt.Println("list:", list)

	hookID, _ := github.CreateWebHook(githubSource, "Sorinttest")
	fmt.Println("hookID:", hookID)*/

	/*exists := github.CheckOrganizationExists(githubSource, "Sorinttest")
	fmt.Println("organization exists:", exists)*/

	/*teams, err := github.GetOrganizationTeams(githubSource, "Sorinttest")
	fmt.Println("err:", err)
	fmt.Println("teams:", teams)*/

	/*users, _ := github.GetTeamMembers(githubSource, 0)
	fmt.Println("users:", users)*/

	commit, err := github.GetCommitMetadata(githubSource, "Sorinttest", "prova1", "113b446f69d7a4c539015827d4fb9f911d0d9eb2")
	fmt.Println("err:", err)
	fmt.Println("commit:", commit)
}

/*func testSynkTrygritea() {
	db := repository.NewAppDb(config.Config)

	gitSource, _ := db.GetGitSourceByName("trygitea")
	if gitSource == nil {
		gitSource = &model.GitSource{Name: "trygitea", GitType: model.Gitea, GitToken: "146f63814f8a52130591d0c5c08d382306e09b48", GitAPIURL: "https://try.gitea.io", AgolaRemoteSource: "alessandropinnawecode", AgolaToken: "79258fee02d5ae315d5dc7eb084e0e20e550cca8"}
		db.SaveGitSource(gitSource)
	}

	fmt.Println("gitSource: ", gitSource)

	organization, _ := db.GetOrganizationByName("SorintDev")
	fmt.Println("organization:", organization)

	repositoryManager.SynkGitRepositorys(&db, organization, gitSource)
}*/

func testSynkWecode() {
	db := repository.NewAppDb(config.Config)

	gitSource, _ := db.GetGitSourceByName("test")
	if gitSource == nil {
		gitSource = &model.GitSource{Name: "test", GitType: model.Gitea, GitToken: "7171eee0d8404600d4da7836f81c07df5a43daaf", GitAPIURL: "https://wecode.sorintdev.it", AgolaRemoteSource: "alessandropinnawecode", AgolaToken: "5adacf4e054f04133c60ac7a4f139dc9aa85b232"}
		db.SaveGitSource(gitSource)
	}

	projects := make(map[string]model.Project)
	projects["agola-example"] = model.Project{GitRepoPath: "agola-example", AgolaProjectID: "c3ff54de-65bb-4a85-8060-ba36cf00610a", Archivied: false}
	projects["Repo1"] = model.Project{GitRepoPath: "Repo1", AgolaProjectID: "c3ff54de-65bb-4a85-8060-ba36cf00610a", Archivied: false}

	/*organization := model.Organization{Name: "TestDemo", GitSourceID: gitSource.ID, Projects: projects, BehaviourType: dto.None, Visibility: dto.Public}
	db.SaveOrganization(&organization)*/

	organization, _ := db.GetOrganizationByName("TestDemo")
	/*organization.Projects = projects
	db.SaveOrganization(organization)*/

	repositoryManager.SynkGitRepositorys(&db, organization, gitSource)
	time.Sleep(10 * time.Minute)
}

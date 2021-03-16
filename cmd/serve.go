package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	giteaApi "wecode.sorint.it/opensource/papagaio-api/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/manager"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/service"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run papagaio server",
	Long:  "run papagaio server",
	Run:   serve,
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	config.SetupConfig()

	if _, err := os.Stat(config.Config.Database.DbPath); os.IsNotExist(err) {
		err := os.Mkdir(config.Config.Database.DbPath, os.ModeDir)
		if err != err {
			panic("Error during mkdir" + config.Config.Database.DbPath)
		}
	}

	db := repository.NewAppDb(config.Config)

	//testGitHub()

	ctrlOrganization := service.OrganizationService{
		Db: &db,
	}

	ctrlGitSource := service.GitSourceService{
		Db: &db,
	}

	ctrlMember := service.MemberService{}

	ctrlWebHook := service.WebHookService{
		Db: &db,
	}

	router := mux.NewRouter()

	controller.SetupHTTPClient()
	controller.SetupRouter(router, &ctrlOrganization, &ctrlGitSource, &ctrlMember, &ctrlWebHook)

	log.Println("Papagaio Server Starting on port ", config.Config.Server.Port)

	logRouter := http.Handler(router)

	if config.Config.LogHTTPRequest {
		logRouter = handlers.LoggingHandler(os.Stdout, router)
	} else {
		logRouter = router
	}

	manager.StartSyncMembers(&db)

	if e := http.ListenAndServe(":"+config.Config.Server.Port, cors.AllowAll().Handler(logRouter)); e != nil {
		log.Println("http server error:", e)
	}

	defer db.DB.Close()
}

func testSomeAPI() {
	gitSource := &model.GitSource{ID: "MTctOTIyZS00ZjFmLWFlODctMDA0N2Q0YTY2MWJk", Name: "gitSourceProva", GitType: "gitea", GitAPIURL: "https://wecode.sorintdev.it", GitToken: "d5e630f316de7132d4f840c305853865b2470cf2"}
	teams, _ := giteaApi.GetOrganizationTeams(gitSource, "Sorint")
	fmt.Println("teams: ", teams)

	for _, team := range *teams {
		users, _ := giteaApi.GetTeamMembers(gitSource, team.ID)
		fmt.Println("Users team: ", users)
	}
}

func testGitHub() {
	githubSource := &model.GitSource{GitToken: "b5433ca882222e373d1cdf37a3406f2ee0423e87"}

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

	users, _ := github.GetOrganizationMembers(githubSource, "Sorinttest")
	fmt.Println("users:", users)
}

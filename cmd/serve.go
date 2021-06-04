package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-api/api/agola"
	"wecode.sorint.it/opensource/papagaio-api/api/git"
	"wecode.sorint.it/opensource/papagaio-api/api/git/gitea"
	"wecode.sorint.it/opensource/papagaio-api/api/git/github"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/controller"
	"wecode.sorint.it/opensource/papagaio-api/repository"
	"wecode.sorint.it/opensource/papagaio-api/service"
	"wecode.sorint.it/opensource/papagaio-api/trigger"
	"wecode.sorint.it/opensource/papagaio-api/utils"
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
	if _, err := os.Stat(config.Config.Database.DbPath); os.IsNotExist(err) {
		err := os.Mkdir(config.Config.Database.DbPath, os.ModeDir)
		if err != err {
			panic("Error during mkdir" + config.Config.Database.DbPath)
		}
	}

	db := repository.NewAppDb(config.Config)
	tr := utils.ConfigUtils{Db: &db}
	agolaApi := agola.AgolaApi{Db: &db}
	gitGateway := git.GitGateway{GiteaApi: &gitea.GiteaApi{Db: &db}, GithubApi: &github.GithubApi{Db: &db}}

	commonMutex := utils.NewEventMutex()

	ctrlOrganization := service.OrganizationService{
		Db:          &db,
		CommonMutex: &commonMutex,
		AgolaApi:    &agolaApi,
		GitGateway:  &gitGateway,
	}

	ctrlGitSource := service.GitSourceService{
		Db:         &db,
		GitGateway: &gitGateway,
	}

	ctrlWebHook := service.WebHookService{
		Db:          &db,
		CommonMutex: &commonMutex,
		AgolaApi:    &agolaApi,
		GitGateway:  &gitGateway,
	}

	ctrlTrigger := service.TriggersService{
		Db: &db,
		Tr: tr,
	}

	sd, err := config.InitTokenSigninData(&config.Config.TokenSigning)
	if err != nil {
		panic(err)
	}
	ctrlOauth2 := service.Oauth2Service{
		Db: &db,
		Sd: sd,
	}

	router := mux.NewRouter()

	controller.SetupHTTPClient()
	controller.SetupRouter(sd, &db, router, &ctrlOrganization, &ctrlGitSource, &ctrlWebHook, &ctrlTrigger, &ctrlOauth2)

	log.Println("Papagaio Server Starting on port ", config.Config.Server.Port)

	logRouter := http.Handler(router)

	if config.Config.LogHTTPRequest {
		logRouter = handlers.LoggingHandler(os.Stdout, router)
	} else {
		logRouter = router
	}

	trigger.StartOrganizationSync(&db, tr, &commonMutex, &agolaApi, &gitGateway)
	trigger.StartRunFailsDiscovery(&db, tr, &commonMutex, &agolaApi, &gitGateway)

	if e := http.ListenAndServe(":"+config.Config.Server.Port, cors.AllowAll().Handler(logRouter)); e != nil {
		log.Println("http server error:", e)
	}

	defer db.DB.Close()
}

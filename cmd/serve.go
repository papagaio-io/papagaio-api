package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-be/config"
	"wecode.sorint.it/opensource/papagaio-be/controller"
	"wecode.sorint.it/opensource/papagaio-be/repository"
	"wecode.sorint.it/opensource/papagaio-be/service"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run papagaio server",
	Long:  "run papagaio server",
	Run:   serve,
}

func Init() {
	rootCmd.AddCommand(serveCmd)
}

func serve(cmd *cobra.Command, args []string) {
	config.Config.Server.Port = "8080" //TODO remove after complited SetupConfig
	//config.SetupConfig()
	db := repository.NewAppDb(config.Config)

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

	if e := http.ListenAndServe(":"+config.Config.Server.Port, cors.AllowAll().Handler(logRouter)); e != nil {
		log.Println("http server error:", e)
	}

	defer db.DB.Close()
}

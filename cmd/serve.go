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
	//Insert only for test
	config.Config.Server.Port = "8080"
	config.Config.Agola.AgolaAddr = "https://agola.sorintdev.it"
	config.Config.Agola.AdminToken = "token admintoken"

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

	testSomeAPI()

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

func testSomeAPI() {
	/*token, err := agolaApi.CreateUserToken("test", "abc")
	fmt.Println("token created for test user: ", token, err)*/

	// gitSource := model.GitSource{GitType: "gitea", GitAPIURL: "https://try.gitea.io", GitToken: "20b93c349872f2bdb3a77b0bd898a3be424c6cbd"}
	// id, _ := gitApi.CreateWebHook(&gitSource, "papagaiotest")
	// fmt.Println("webhook id: ", id)

	// idOrganization, err := agolaApi.CreateOrganization("DC-Comics", "public")
	// fmt.Println("Create Organization: id ", idOrganization, err)

	// err1 := agolaApi.AddOrganizationMember("TullioOrg2", "tullio", "member")
	// fmt.Println("Add org member Error ", err1)

	// err1 := agolaApi.RemoveOrganizationMember("TullioOrg2", "tullio")
	// fmt.Println("Add org member Error ", err1)
}

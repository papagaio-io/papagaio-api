package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-be/controller"
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
	fmt.Println("Papagaio Server Starting")
	router := controller.NewRouter()
	if err := http.ListenAndServe(":9000", router); err != nil {
		panic(err)
	}

}

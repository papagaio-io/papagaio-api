package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-api/config"
)

var rootCmd = &cobra.Command{
	Use:     "papagaio",
	Short:   "Papagaio",
	Long:    "Papagaio Extentions Agola",
	Version: "1",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	config.SetupConfig()
}

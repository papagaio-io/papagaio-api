package cmd

import (
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

var gitSourceCmd = &cobra.Command{
	Use: "gitsource",
}

var addUGitSourceCmd = &cobra.Command{
	Use: "add",
	Run: addGitSource,
}

var removeGitSourceCmd = &cobra.Command{
	Use: "remove",
	Run: removeGitSource,
}

var cfgGitSource configGitSourceCmd

type configGitSourceCmd struct {
	name              string
	gitType           string
	gitAPIURL         string
	gitToken          string
	agolaRemoteSource string
	agolaToken        string
}

func init() {
	rootCmd.AddCommand(gitSourceCmd)
	gitSourceCmd.AddCommand(addUGitSourceCmd)
	gitSourceCmd.AddCommand(removeGitSourceCmd)

	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.name, "name", "", "gitSource name")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.gitType, "type", "", "git type")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.gitAPIURL, "git-api-url", "", "api url")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.gitToken, "git-token", "", "git token")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.agolaRemoteSource, "agola-remotesource", "", "Agola remotesource")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.agolaToken, "agola-token", "", "Agola token")
}

func addGitSource(cmd *cobra.Command, args []string) {
	if len(cfgGitSource.name) == 0 {
		cmd.PrintErrln("name is empty or not valid")
		os.Exit(1)
	}

	if strings.Compare(cfgGitSource.gitType, "gitea") != 0 && strings.Compare(cfgGitSource.gitType, "github") != 0 {
		cmd.PrintErrln("type must be gitea or github")
		os.Exit(1)
	}

	if len(cfgGitSource.gitAPIURL) > 0 {
		_, err := url.ParseRequestURI(cfgGitSource.gitAPIURL)
		if err != nil {
			cmd.PrintErrln("git-api-url is not valid")
			os.Exit(1)
		}
	} else {
		if strings.Compare(cfgGitSource.gitType, "gitea") == 0 {
			cmd.PrintErrln("gitea type need git-api-url")
			os.Exit(1)
		}
	}

	if len(cfgGitSource.gitToken) == 0 {
		cmd.PrintErrln("git-token not valid")
		os.Exit(1)
	}

	if len(cfgGitSource.agolaRemoteSource) == 0 {
		cmd.PrintErrln("agola-remotesource not valid")
		os.Exit(1)
	}

	if len(cfgGitSource.agolaToken) == 0 {
		cmd.PrintErrln("agola-token not valid")
		os.Exit(1)
	}

	db := beginGitSource(cmd)

	gitSource, _ := db.GetGitSourceByName(cfgGitSource.name)
	if gitSource == nil {
		db.SaveGitSource(&model.GitSource{
			Name:              cfgGitSource.name,
			GitType:           model.GitType(cfgGitSource.gitType),
			GitAPIURL:         cfgGitSource.gitAPIURL,
			GitToken:          cfgGitSource.gitToken,
			AgolaRemoteSource: cfgGitSource.agolaRemoteSource,
			AgolaToken:        cfgGitSource.agolaToken,
		})
		cmd.Println("GitSource ", cfgGitSource.name, "saved")
	} else {
		cmd.PrintErrln("GitSource", cfgGitSource.name, "just present in db")
	}
}

func removeGitSource(cmd *cobra.Command, args []string) {
	if len(cfgGitSource.name) == 0 {
		cmd.PrintErrln("name is empty or not valid")
		os.Exit(1)
	}

	db := beginGitSource(cmd)

	gitSource, _ := db.GetGitSourceByName(cfgGitSource.name)
	if gitSource != nil {
		db.DeleteGitSource(gitSource.ID)
		cmd.Println("GitSource ", cfgGitSource.name, "removed")
	} else {
		cmd.PrintErrln("GitSource", cfgGitSource.name, "not found")
	}
}

func beginGitSource(cmd *cobra.Command) repository.AppDb {
	config.SetupConfig()

	if _, err := os.Stat(config.Config.Database.DbPath); os.IsNotExist(err) {
		err := os.Mkdir(config.Config.Database.DbPath, os.ModeDir)
		if err != err {
			panic("Error during mkdir" + config.Config.Database.DbPath)
		}
	}

	return repository.NewAppDb(config.Config)
}

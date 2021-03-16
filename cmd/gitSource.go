package cmd

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-api/config"
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
	userCmd.AddCommand(addUGitSourceCmd)
	userCmd.AddCommand(removeGitSourceCmd)

	userCmd.PersistentFlags().StringVar(&cfgGitSource.name, "name", "", "gitSource name")
	userCmd.PersistentFlags().StringVar(&cfgGitSource.gitType, "type", "", "git type")
	userCmd.PersistentFlags().StringVar(&cfgGitSource.gitAPIURL, "git-api-url", "", "api url")
	userCmd.PersistentFlags().StringVar(&cfgGitSource.gitToken, "git-token", "", "git token")
	userCmd.PersistentFlags().StringVar(&cfgGitSource.agolaRemoteSource, "agola-remotesource", "", "Agola remotesource")
	userCmd.PersistentFlags().StringVar(&cfgGitSource.agolaToken, "agola-token", "", "Agola token")
}

func addGitSource(cmd *cobra.Command, args []string) {
	beginGitSource(cmd)
}

func removeGitSource(cmd *cobra.Command, args []string) {
	beginGitSource(cmd)
}

var regexUrl = regexp.MustCompile(`/((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)/`)

func beginGitSource(cmd *cobra.Command) repository.AppDb {
	if len(cfgGitSource.name) == 0 {
		cmd.PrintErrln("name is empty or not valid")
	}

	if strings.Compare(cfgGitSource.gitType, "gitea") != 0 || strings.Compare(cfgGitSource.gitType, "github") != 0 {
		cmd.PrintErrln("type must be gitea or github")
	}

	if !emailRegex.MatchString(cfgGitSource.gitAPIURL) {
		cmd.PrintErrln("egit-api-url is not valid")
	}

	config.SetupConfig()

	if _, err := os.Stat(config.Config.Database.DbPath); os.IsNotExist(err) {
		log.Println("Filder", config.Config.Database.DbPath, "not found")
		panic("Filder" + config.Config.Database.DbPath + "not found")
	}

	return repository.NewAppDb(config.Config)
}

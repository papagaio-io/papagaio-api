package cmd

import "github.com/spf13/cobra"

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
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(addUserCmd)
	userCmd.AddCommand(removeUserCmd)

	userCmd.PersistentFlags().StringVar(&cfgGitSource.name, "name", "", "gitSource name")
	userCmd.PersistentFlags().StringVar(&cfgGitSource.gitType, "type", "", "git type")
	userCmd.PersistentFlags().StringVar(&cfgGitSource.gitAPIURL, "git-api-url", "", "api url")
	userCmd.PersistentFlags().StringVar(&cfgGitSource.gitToken, "git-token", "", "git token")
	userCmd.PersistentFlags().StringVar(&cfgGitSource.agolaRemoteSource, "agola-remotesource", "", "Agola remotesource")
	userCmd.PersistentFlags().StringVar(&cfgGitSource.agolaToken, "agola-token", "", "Agola token")
}

func addGitSource(cmd *cobra.Command, args []string) {

}

func removeGitSource(cmd *cobra.Command, args []string) {

}

/*func beginGitSource(cmd *cobra.Command) repository.AppDb {
	if !emailRegex.MatchString(cfgUser.email) {
		cmd.PrintErrln("email is empty or not valid")
	}

	config.SetupConfig()

	if _, err := os.Stat(config.Config.Database.DbPath); os.IsNotExist(err) {
		log.Println("Filder", config.Config.Database.DbPath, "not found")
		panic("Filder" + config.Config.Database.DbPath + "not found")
	}

	return repository.NewAppDb(config.Config)
}*/

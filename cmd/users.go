package cmd

import (
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/model"
	"wecode.sorint.it/opensource/papagaio-api/repository"
)

var userCmd = &cobra.Command{
	Use: "user",
}

var addUserCmd = &cobra.Command{
	Use: "add",
	Run: addUser,
}

var removeUserCmd = &cobra.Command{
	Use: "remove",
	Run: removeUser,
}

var cfgUser configUserCmd
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type configUserCmd struct {
	email string
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(addUserCmd)
	userCmd.AddCommand(removeUserCmd)

	userCmd.PersistentFlags().StringVar(&cfgUser.email, "email", "", "user email")
}

func addUser(cmd *cobra.Command, args []string) {
	db := beginUser(cmd)
	user, _ := db.GetUserByEmail(cfgUser.email)
	if user == nil {
		db.SaveUser(&model.User{Email: cfgUser.email})
		cmd.Println("User ", cfgUser.email, "saved")
	} else {
		cmd.PrintErrln("User", cfgUser.email, "just present in db")
	}
}

func removeUser(cmd *cobra.Command, args []string) {
	db := beginUser(cmd)

	user, _ := db.GetUserByEmail(cfgUser.email)
	if user != nil {
		db.DeleteUser(cfgUser.email)
		cmd.Println("User ", cfgUser.email, "removed")
	} else {
		cmd.PrintErrln("User", cfgUser.email, "not found")
	}
}

func beginUser(cmd *cobra.Command) repository.AppDb {
	if !emailRegex.MatchString(cfgUser.email) {
		cmd.PrintErrln("email is empty or not valid")
		os.Exit(1)
	}

	config.SetupConfig()

	if _, err := os.Stat(config.Config.Database.DbPath); os.IsNotExist(err) {
		cmd.PrintErrln("Filder", config.Config.Database.DbPath, "not found")
		os.Exit(1)
	}

	return repository.NewAppDb(config.Config)
}

package cmd

import (
	"log"
	"os"
	"regexp"

	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-be/config"
	"wecode.sorint.it/opensource/papagaio-be/model"
	"wecode.sorint.it/opensource/papagaio-be/repository"
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

var cfg configUserCmd
var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type configUserCmd struct {
	email string
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(addUserCmd)
	userCmd.AddCommand(removeUserCmd)

	userCmd.PersistentFlags().StringVar(&cfg.email, "email", "", "user email")
}

func addUser(cmd *cobra.Command, args []string) {
	db := begin(cmd)
	user, _ := db.GetUserByEmail(cfg.email)
	if user == nil {
		db.SaveUser(&model.User{Email: cfg.email})
		cmd.Println("User ", cfg.email, "saved")
	} else {
		cmd.PrintErrln("User", cfg.email, "just present in db")
	}
}

func removeUser(cmd *cobra.Command, args []string) {
	db := begin(cmd)

	user, _ := db.GetUserByEmail(cfg.email)
	if user != nil {
		db.DeleteUser(cfg.email)
		cmd.Println("User ", cfg.email, "removed")
	} else {
		cmd.PrintErrln("User", cfg.email, "not found")
	}
}

func begin(cmd *cobra.Command) repository.AppDb {
	if !emailRegex.MatchString(cfg.email) {
		cmd.PrintErrln("email is empty or not valid")
	}

	config.SetupConfig()

	if _, err := os.Stat(config.Config.Database.DbPath); os.IsNotExist(err) {
		log.Println("Filder", config.Config.Database.DbPath, "not found")
		panic("Filder" + config.Config.Database.DbPath + "not found")
	}

	return repository.NewAppDb(config.Config)
}

package cmd

import (
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-api/api"
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

const constEmailRegex = "^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

var emailRegex = regexp.MustCompile(constEmailRegex)

type configUserCmd struct {
	CommonConfig

	email string
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(addUserCmd)
	userCmd.AddCommand(removeUserCmd)

	AddCommonFlags(userCmd, &cfgUser.CommonConfig)
	userCmd.PersistentFlags().StringVar(&cfgUser.email, "email", "", "user email")
}

func addUser(cmd *cobra.Command, args []string) {
	beginUser(cmd)

	client := &http.Client{}
	URLApi := cfgUser.gatewayURL + "/api/adduser"
	reqBody := strings.NewReader(`{"email": "` + cfgUser.email + `"}`)
	req, _ := http.NewRequest("POST", URLApi, reqBody)
	req.Header.Add("Authorization", "token "+cfgUser.token)

	resp, _ := client.Do(req)

	if !api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		cmd.PrintErrln(string(body))
	} else {
		cmd.Println("user created")
	}
}

func removeUser(cmd *cobra.Command, args []string) {
	beginUser(cmd)

	client := &http.Client{}
	URLApi := cfgUser.gatewayURL + "/api/removeuser/" + cfgUser.email
	req, _ := http.NewRequest("DELETE", URLApi, nil)
	req.Header.Add("Authorization", "token "+cfgUser.token)

	resp, _ := client.Do(req)
	if !api.IsResponseOK(resp.StatusCode) {
		body, _ := ioutil.ReadAll(resp.Body)
		cmd.PrintErrln(string(body))
	} else {
		cmd.Println("user removed")
	}
}

func beginUser(cmd *cobra.Command) {

	if err := cfgUser.IsAdminUser(); err != nil {
		cmd.PrintErrln(err.Error())
		os.Exit(1)
	}

	if !emailRegex.MatchString(cfgUser.email) {
		cmd.PrintErrln("email is empty or not valid")
		os.Exit(1)
	}
}

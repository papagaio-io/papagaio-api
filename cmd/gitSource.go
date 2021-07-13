package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-api/api"
	"wecode.sorint.it/opensource/papagaio-api/config"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/types"
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

var updateGitSourceCmd = &cobra.Command{
	Use: "update",
	Run: updateGitSource,
}

var cfgGitSource configGitSourceCmd

type configGitSourceCmd struct {
	CommonConfig

	name                  string
	gitType               string
	gitAPIURL             string
	gitClientID           string
	gitClientSecret       string
	agolaRemoteSourceName string
	agolaClientID         string
	agolaClientSecret     string

	deleteRemoteSource bool
}

func init() {
	config.SetupConfig()

	rootCmd.AddCommand(gitSourceCmd)
	gitSourceCmd.AddCommand(addUGitSourceCmd)
	gitSourceCmd.AddCommand(removeGitSourceCmd)
	gitSourceCmd.AddCommand(updateGitSourceCmd)

	AddCommonFlags(gitSourceCmd, &cfgGitSource.CommonConfig)

	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.name, "name", "", "gitSource name")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.gitType, "type", "", "git type(gitea, github)")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.gitAPIURL, "git-api-url", "", "api url")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.gitClientID, "git-client-id", "", "git oauth2 client id")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.gitClientSecret, "git-client-secret", "", "git oauth2 client secret")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.agolaRemoteSourceName, "agola-remotesource", "", "agola remotesource name")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.agolaClientID, "agola-client-id", "", "agola oauth2 client id")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.agolaClientSecret, "agola-client-secret", "", "agola oauth2 client secret")

	gitSourceCmd.PersistentFlags().BoolVar(&cfgGitSource.deleteRemoteSource, "delete-remotesource", false, "true to delete the Agola remotesource")
}

func addGitSource(cmd *cobra.Command, args []string) {
	if err := cfgGitSource.IsAdminUser(); err != nil {
		cmd.PrintErrln(err.Error())
		os.Exit(1)
	}

	gitSourceRequest := dto.CreateGitSourceRequestDto{
		Name:                  cfgGitSource.name,
		GitType:               types.GitType(cfgGitSource.gitType),
		GitAPIURL:             &cfgGitSource.gitAPIURL,
		GitClientID:           cfgGitSource.gitClientID,
		GitClientSecret:       cfgGitSource.gitClientSecret,
		AgolaRemoteSourceName: &cfgGitSource.agolaRemoteSourceName,
		AgolaClientID:         &cfgGitSource.agolaClientID,
		AgolaClientSecret:     &cfgGitSource.agolaClientSecret,
	}

	err := gitSourceRequest.IsValid()
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	data, _ := json.Marshal(gitSourceRequest)

	client := &http.Client{}
	URLApi := cfgGitSource.gatewayURL + "/api/gitsource"
	reqBody := strings.NewReader(string(data))
	req, _ := http.NewRequest("POST", URLApi, reqBody)
	req.Header.Add("Authorization", "token "+cfgGitSource.token)

	resp, err := client.Do(req)
	if err != nil {
		cmd.Println("Error:", err.Error())
	} else {
		if !api.IsResponseOK(resp.StatusCode) {
			body, _ := ioutil.ReadAll(resp.Body)
			cmd.PrintErrln("Somefing was wrong! " + string(body))
			os.Exit(1)
		}

		cmd.Println("gitsource created")
	}
}

func removeGitSource(cmd *cobra.Command, args []string) {
	if err := cfgGitSource.IsAdminUser(); err != nil {
		cmd.PrintErrln(err.Error())
		os.Exit(1)
	}

	if len(cfgGitSource.name) == 0 {
		cmd.PrintErrln("name is empty or not valid")
		os.Exit(1)
	}

	deleteRemoteSourceParam := ""
	if cfgGitSource.deleteRemoteSource {
		deleteRemoteSourceParam = "?deleteremotesource"
	}

	client := &http.Client{}
	URLApi := cfgGitSource.gatewayURL + "/api/gitsource/" + cfgGitSource.name + deleteRemoteSourceParam
	req, _ := http.NewRequest("DELETE", URLApi, nil)
	req.Header.Add("Authorization", "token "+cfgGitSource.token)

	resp, err := client.Do(req)
	if err != nil {
		cmd.Println("Error:", err.Error())
	} else {
		if !api.IsResponseOK(resp.StatusCode) {
			body, _ := ioutil.ReadAll(resp.Body)
			cmd.PrintErrln("Something was wrong! " + string(body))
			os.Exit(1)
		}
	}

	cmd.Println("gitsource removed")
}

func updateGitSource(cmd *cobra.Command, args []string) {
	if err := cfgGitSource.IsAdminUser(); err != nil {
		cmd.PrintErrln(err.Error())
		os.Exit(1)
	}

	if len(cfgGitSource.name) == 0 {
		cmd.PrintErrln("name is empty or not valid")
		os.Exit(1)
	}

	requestDto := dto.UpdateGitSourceRequestDto{}

	if len(cfgGitSource.gitAPIURL) != 0 {
		requestDto.GitAPIURL = &cfgGitSource.gitAPIURL
	}
	if len(cfgGitSource.gitType) != 0 {
		if strings.Compare(cfgGitSource.gitType, "gitea") != 0 && strings.Compare(cfgGitSource.gitType, "github") != 0 {
			cmd.PrintErrln("type must be gitea or github")
			os.Exit(1)
		}
		requestDto.GitType = (*types.GitType)(&cfgGitSource.gitType)
	}
	if len(cfgGitSource.gitClientID) != 0 {
		requestDto.GitClientID = &cfgGitSource.gitClientID
	}
	if len(cfgGitSource.gitClientSecret) != 0 {
		requestDto.GitClientSecret = &cfgGitSource.gitClientSecret
	}

	data, _ := json.Marshal(requestDto)

	client := &http.Client{}
	URLApi := cfgGitSource.gatewayURL + "/api/gitsource/" + cfgGitSource.name
	reqBody := strings.NewReader(string(data))
	req, _ := http.NewRequest("PUT", URLApi, reqBody)
	req.Header.Add("Authorization", "token "+cfgGitSource.token)

	resp, err := client.Do(req)
	if err != nil {
		cmd.Println("Error:", err.Error())
	} else {
		if !api.IsResponseOK(resp.StatusCode) {
			body, _ := ioutil.ReadAll(resp.Body)
			cmd.PrintErrln("Somefing was wrong! " + string(body))
			os.Exit(1)
		}
	}

	cmd.Println("gitsource updated")
}

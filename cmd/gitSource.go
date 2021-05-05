package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-api/api"
	"wecode.sorint.it/opensource/papagaio-api/dto"
	"wecode.sorint.it/opensource/papagaio-api/model"
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
	gitSourceCmd.AddCommand(updateGitSourceCmd)

	AddCommonFlags(gitSourceCmd, &cfgGitSource.CommonConfig)

	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.name, "name", "", "gitSource name")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.gitType, "type", "", "git type")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.gitAPIURL, "git-api-url", "", "api url")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.gitToken, "git-token", "", "git token")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.agolaRemoteSource, "agola-remotesource", "", "Agola remotesource")
	gitSourceCmd.PersistentFlags().StringVar(&cfgGitSource.agolaToken, "agola-token", "", "Agola token")
}

func addGitSource(cmd *cobra.Command, args []string) {
	if err := cfgGitSource.IsAdminUser(); err != nil {
		cmd.PrintErrln(err.Error())
		os.Exit(1)
	}

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

	gitSource := model.GitSource{
		Name:              cfgGitSource.name,
		GitType:           types.GitType(cfgGitSource.gitType),
		GitAPIURL:         cfgGitSource.gitAPIURL,
		GitToken:          cfgGitSource.gitToken,
		AgolaRemoteSource: cfgGitSource.agolaRemoteSource,
		AgolaToken:        cfgGitSource.agolaToken,
	}
	data, _ := json.Marshal(gitSource)

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

	client := &http.Client{}
	URLApi := cfgGitSource.gatewayURL + "/api/gitsource/" + cfgGitSource.name
	req, _ := http.NewRequest("DELETE", URLApi, nil)
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

	cmd.Println("gitsource removed")
}

//TODO
func updateGitSource(cmd *cobra.Command, args []string) {
	if err := cfgGitSource.IsAdminUser(); err != nil {
		cmd.PrintErrln(err.Error())
		os.Exit(1)
	}

	if len(cfgGitSource.name) == 0 {
		cmd.PrintErrln("name is empty or not valid")
		os.Exit(1)
	}

	requestDto := dto.UpdateRemoteSourceRequestDto{}

	if len(cfgGitSource.gitAPIURL) != 0 {
		requestDto.GitAPIURL = &cfgGitSource.gitAPIURL
	}
	if len(cfgGitSource.gitToken) != 0 {
		requestDto.GitToken = &cfgGitSource.gitToken
	}
	if len(cfgGitSource.gitType) != 0 {
		if strings.Compare(cfgGitSource.gitType, "gitea") != 0 && strings.Compare(cfgGitSource.gitType, "github") != 0 {
			cmd.PrintErrln("type must be gitea or github")
			os.Exit(1)
		}
		requestDto.GitType = (*types.GitType)(&cfgGitSource.gitToken)
	}
	if len(cfgGitSource.agolaRemoteSource) != 0 {
		requestDto.AgolaRemoteSource = &cfgGitSource.agolaRemoteSource
	}
	if len(cfgGitSource.agolaToken) != 0 {
		requestDto.AgolaToken = &cfgGitSource.agolaToken
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

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
)

var userCmd = &cobra.Command{
	Use: "user",
}

var changeUserRolCmd = &cobra.Command{
	Use: "change-role",
	Run: changeUserRole,
}

var cfgUser configUser

type configUser struct {
	CommonConfig

	userId   uint64
	userRole string
}

func init() {
	config.SetupConfig()

	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(changeUserRolCmd)

	AddCommonFlags(userCmd, &cfgUser.CommonConfig)

	userCmd.PersistentFlags().Uint64Var(&cfgUser.userId, "id", uint64(0), "user id")
	userCmd.PersistentFlags().StringVar(&cfgUser.userRole, "role", "", "user role(ADMINISTRATOR, DEVELOPER)")
}

func changeUserRole(cmd *cobra.Command, args []string) {
	if err := cfgGitSource.IsAdminUser(); err != nil {
		cmd.PrintErrln(err.Error())
		os.Exit(1)
	}

	request := dto.ChangeUserRoleRequestDto{
		UserID:   &cfgUser.userId,
		UserRole: dto.UserRole(cfgUser.userRole),
	}

	err := request.IsValid()
	if err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}

	data, _ := json.Marshal(request)
	client := &http.Client{}
	URLApi := cfgGitSource.gatewayURL + "/api/changeuserrole"
	reqBody := strings.NewReader(string(data))
	req, _ := http.NewRequest("PUT", URLApi, reqBody)
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

		cmd.Println("user role changed")
	}
}

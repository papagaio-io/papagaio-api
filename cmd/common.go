package cmd

import (
	"errors"
	"strings"

	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-api/config"
)

type CommonConfig struct {
	gatewayURL string
	token      string
}

func AddCommonFlags(cmd *cobra.Command, cfg *CommonConfig) {
	cmd.PersistentFlags().StringVar(&cfg.gatewayURL, "gateway-url", config.Config.CmdConfig.DefaultGatewayURL, "papagaio gateway URL")
	cmd.PersistentFlags().StringVar(&cfg.token, "token", "", "token")
}

func (common CommonConfig) IsAdminUser() error {
	if len(common.token) == 0 {
		return errors.New("token is required")
	}

	if strings.Compare(common.token, config.Config.CmdConfig.Token) != 0 {
		return errors.New("token not valit! must be an admin user")
	}
	return nil
}

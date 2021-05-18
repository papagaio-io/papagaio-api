package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"wecode.sorint.it/opensource/papagaio-api/config"
)

type CommonConfig struct {
	gatewayURL string
	token      string
}

func AddCommonFlags(cmd *cobra.Command, cfg *CommonConfig) {
	cmd.PersistentFlags().StringVar(&cfg.gatewayURL, "gateway-url", config.Config.CmdConfig.DefaultGatewayURL, "papagaio gateway URL(optional)")
	cmd.PersistentFlags().StringVar(&cfg.token, "token", "", "token")
}

func (common CommonConfig) IsAdminUser() error {
	if len(common.token) == 0 {
		return errors.New("token is required")
	}

	return nil
}

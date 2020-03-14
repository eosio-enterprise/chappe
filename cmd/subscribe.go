package cmd

import (
	"context"
	"os"

	"github.com/eosio-enterprise/chappe/pkg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// MakeSubscribe ...
func MakeSubscribe() *cobra.Command {
	var command = &cobra.Command{
		Use:          "subscribe",
		Short:        "Subscribe to a channel",
		Example:      ` chappe subscribe --channel-name MyChannel`,
		SilenceUsage: false,
	}

	var channelName string
	command.Flags().StringVarP(&channelName, "channel-name", "n", "", "channel name")
	command.Run = func(cmd *cobra.Command, args []string) {

		go func() {

			if viper.GetString("Dfuse.Protocol") == "WebSocket" {
				pkg.StreamWS(channelName)
			} else {
				pkg.StreamMessages(context.TODO(), channelName)
			}
		}()

		sigs := make(chan os.Signal, 1)
		<-sigs
	}
	return command
}

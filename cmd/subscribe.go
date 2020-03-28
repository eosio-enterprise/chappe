package cmd

import (
	"context"
	"os"

	"github.com/eosio-enterprise/chappe/pkg"
	"github.com/spf13/cobra"
)

// MakeSubscribe ...
func MakeSubscribe() *cobra.Command {
	var command = &cobra.Command{
		Use:          "subscribe",
		Short:        "Subscribe to a channel",
		Example:      ` chappe subscribe --channel-name MyChannel --send-receipts --ledgerFile my_ledger_file.dat`,
		SilenceUsage: false,
	}

	var channelName, ledgerFile string
	var sendReceipts bool
	command.Flags().StringVarP(&channelName, "channel-name", "n", "", "channel name")
	command.Flags().BoolVarP(&sendReceipts, "send-receipts", "r", false, "send device-specific receipts back to hub")
	command.Flags().StringVarP(&ledgerFile, "ledger-file", "l", "", "record received transactions to a ledger file")

	command.Run = func(cmd *cobra.Command, args []string) {

		go func() {
			pkg.StreamMessages(context.TODO(), channelName, ledgerFile, sendReceipts)
		}()

		sigs := make(chan os.Signal, 1)
		<-sigs
	}
	return command
}

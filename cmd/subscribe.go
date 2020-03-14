package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dfuse-io/eosws-go"
	"github.com/eosio-enterprise/chappe/pkg"
	"github.com/spf13/cobra"
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

			pkg.StreamMessages(context.TODO())

			client := pkg.GetClient()
			err := client.Send(pkg.GetActionTraces())
			if err != nil {
				log.Fatalf("Failed to send request to dfuse: %s", err)
			}

			for {
				msg, err := client.Read()
				if err != nil {
					log.Fatalf("Cannot read from dfuse client: %s", err)
				}

				switch m := msg.(type) {
				case *eosws.ActionTrace:
					pkg.StreamMessages(context.TODO())
					// pkg.Receive(channelName, m)
				case *eosws.Progress:
					fmt.Print(".") // poor man's progress bar, using print not log
				case *eosws.Listening:
					log.Println("Received Listening Message ...")
				default:
					log.Println("Received Unsupported Message", m)
				}
			}
		}()

		sigs := make(chan os.Signal, 1)
		<-sigs
	}
	return command
}

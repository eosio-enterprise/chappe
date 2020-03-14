package cmd

import (
	"fmt"
	"log"

	"github.com/eosio-enterprise/chappe/internal/encryption"
	"github.com/spf13/cobra"
)

// MakeCreate ...
func MakeCreate() *cobra.Command {
	var command = &cobra.Command{
		Use:          "create",
		Short:        "Create a new private chappe channel",
		Example:      `  chappe create --channel-name chan649`,
		SilenceUsage: false,
	}

	var channelName string
	command.Flags().StringVarP(&channelName, "channel-name", "n", "", "Name of the private channel to create")

	command.RunE = func(command *cobra.Command, args []string) error {

		if len(channelName) == 0 {
			return fmt.Errorf("--channel-name required")
		}

		encryption.CreateChannel(channelName)

		log.Println(
			`=======================================================================
key ` + channelName + ` created in files ` + channelName + `.pem (private) and ` + channelName + `.pub (public)
=======================================================================`)

		return nil
	}
	return command
}

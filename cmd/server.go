package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MakeServer ...
func MakeServer() *cobra.Command {
	var command = &cobra.Command{
		Use:          "server",
		Short:        "Run a server",
		Example:      `  chappe server`,
		SilenceUsage: false,
	}
	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("Doesn't do anything. TODO: implement subscribed as a chappe command.")
		fmt.Println("Ctrl-C to exit")
		fmt.Println("Contribute: https://github.com/eosio-enterprise/chappe/blob/master/CONTRIBUTING.md")

	}
	return command
}

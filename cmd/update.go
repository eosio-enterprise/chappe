package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MakeUpdate ...
func MakeUpdate() *cobra.Command {
	var command = &cobra.Command{
		Use:          "update",
		Short:        "Print update instructions",
		Example:      `  chappe update`,
		SilenceUsage: false,
	}
	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("Subscribe action is not yet implemented")
		fmt.Println("Contribute: https://github.com/eosio-enterprise/chappe/blob/master/CONTRIBUTING.md")
	}
	return command
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MakeSubscribe ...
func MakeSubscribe() *cobra.Command {
	var command = &cobra.Command{
		Use:          "subscribe",
		Short:        "Subscribe to a private communications channel",
		Example:      `  chappe subscribe`,
		SilenceUsage: false,
	}
	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("Subscribe action is not yet implemented")
	}
	return command
}

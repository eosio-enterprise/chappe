package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// MakePublish ...
func MakePublish() *cobra.Command {
	var command = &cobra.Command{
		Use:          "publish",
		Short:        "Publish a private message to a channel",
		Example:      `  chappe publish`,
		SilenceUsage: false,
	}
	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("Publish action is not yet implemented")
	}
	return command
}

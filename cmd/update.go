package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func MakeUpdate() *cobra.Command {
	var command = &cobra.Command{
		Use:          "update",
		Short:        "Print update instructions",
		Example:      `  chappe update`,
		SilenceUsage: false,
	}
	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println(chappeUpdate)
	}
	return command
}

const chappeUpdate = `Update instructions coming soon`

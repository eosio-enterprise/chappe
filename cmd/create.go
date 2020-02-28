package cmd

import (
	"fmt"
	"strings"

	"github.com/eosio-enterprise/chappe/cmd/apps"
	"github.com/spf13/cobra"
)

func MakeCreate() *cobra.Command {
	var command = &cobra.Command{
		Use:   "create",
		Short: "Create chappe artifacts",
		Long: `Create and share channel certificates with other nodes
on the chappe network. Potentially create new networks.`,
		Example: `  chappe create
  chappe create channel --channel-name MyChannel
  chappe create key --key-name MyKey`,
		SilenceUsage: false,
	}

	command.PersistentFlags().String("kubeconfig", "kubeconfig", "Local path for your kubeconfig file")

	command.RunE = func(command *cobra.Command, args []string) error {

		if len(args) == 0 {
			fmt.Printf("You can create: %s\n%s\n\n", strings.TrimRight("\n - "+strings.Join(getApps(), "\n - "), "\n - "),
				`Run chappe create NAME --help to see configuration options.`)
			return nil
		}

		return nil
	}

	command.AddCommand(apps.MakeCreateKey())
	// command.AddCommand(apps.MakeInstallMetricsServer())

	// command.AddCommand(MakeInfo())

	return command
}

func getApps() []string {
	return []string{"channel",
		"key",
	}
}

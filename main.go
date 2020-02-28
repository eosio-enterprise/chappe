package main

import (
	"os"

	cmd "github.com/eosio-enterprise/chappe/cmd"
	"github.com/spf13/cobra"
)

func main() {

	cmdVersion := cmd.MakeVersion()
	cmdCreate := cmd.MakeCreate()
	// cmdInfo := cmd.MakeInfo()
	cmdUpdate := cmd.MakeUpdate()

	// printarkadeASCIIArt := cmd.PrintArkadeASCIIArt

	var rootCmd = &cobra.Command{
		Use: "chappe",
		Run: func(cmd *cobra.Command, args []string) {
			printarkadeASCIIArt()
			cmd.Help()
		},
	}

	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdVersion)
	// rootCmd.AddCommand(cmdInfo)
	rootCmd.AddCommand(cmdUpdate)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

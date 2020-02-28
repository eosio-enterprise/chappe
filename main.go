package main

import (
	"os"

	cmd "github.com/eosio-enterprise/chappe/cmd"
	"github.com/spf13/cobra"
)

func main() {

	cmdVersion := cmd.MakeVersion()
	cmdCreate := cmd.MakeCreate()
	cmdUpdate := cmd.MakeUpdate()
	cmdPublish := cmd.MakePublish()
	cmdSubscribe := cmd.MakeSubscribe()

	printChappeASCIIArt := cmd.PrintChappeASCIIArt

	var rootCmd = &cobra.Command{
		Use: "chappe",
		Run: func(cmd *cobra.Command, args []string) {
			printChappeASCIIArt()
			cmd.Help()
		},
	}

	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdVersion)
	rootCmd.AddCommand(cmdUpdate)
	rootCmd.AddCommand(cmdPublish)
	rootCmd.AddCommand(cmdSubscribe)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"

	cmd "github.com/eosio-enterprise/chappe/cmd"
	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func main() {

	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("config file not found: %s", err))
	}

	cmdVersion := cmd.MakeVersion()
	cmdCreate := cmd.MakeCreate()
	cmdUpdate := cmd.MakeUpdate()
	cmdPublish := cmd.MakePublish()
	cmdServer := cmd.MakeServer()
	cmdSubscribe := cmd.MakeSubscribe()
	cmdGet := cmd.MakeGet()

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
	rootCmd.AddCommand(cmdServer)
	rootCmd.AddCommand(cmdSubscribe)
	rootCmd.AddCommand(cmdGet)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

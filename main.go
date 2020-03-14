package main

import (
	"log"
	"os"

	cmd "github.com/eosio-enterprise/chappe/cmd"
	"github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func main() {

	cmdVersion := cmd.MakeVersion()
	cmdCreate := cmd.MakeCreate()
	cmdUpdate := cmd.MakeUpdate()
	cmdPublish := cmd.MakePublish()
	cmdSubscribe := cmd.MakeSubscribe()
	cmdGet := cmd.MakeGet()

	var rootCmd = &cobra.Command{
		Use: "chappe",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	var configFile string
	rootCmd.Flags().StringVarP(&configFile, "config-file", "c", "configs/config.yaml", "Path/name of configuration file")
	viper.SetConfigFile(configFile) //"configs/config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("config file not found: %s", err)
	}

	rootCmd.AddCommand(cmdCreate)
	rootCmd.AddCommand(cmdVersion)
	rootCmd.AddCommand(cmdUpdate)
	rootCmd.AddCommand(cmdPublish)
	rootCmd.AddCommand(cmdSubscribe)
	rootCmd.AddCommand(cmdGet)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

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

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("/etc/chappe/")
	viper.AddConfigPath("$HOME/.chappe")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("CHAPPE")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
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

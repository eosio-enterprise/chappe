package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version ...
	Version string
	// GitCommit ...
	GitCommit string
)

// PrintChappeASCIIArt ...
func PrintChappeASCIIArt() {
	// chappeLogo := aec.BlueF.Apply(chappeFigletStr)
	fmt.Print(chappeFigletStr)
}

// MakeVersion ...
func MakeVersion() *cobra.Command {
	var command = &cobra.Command{
		Use:          "version",
		Short:        "Print the version",
		Example:      `  chappe version`,
		SilenceUsage: false,
	}
	command.Run = func(cmd *cobra.Command, args []string) {
		PrintChappeASCIIArt()
		if len(Version) == 0 {
			fmt.Println("Version: dev")
		} else {
			fmt.Println("Version:", Version)
		}
		fmt.Println("Git Commit:", GitCommit)
	}
	return command
}

// TODO: Print chappe figlet logo with version
const chappeFigletStr = "Welcome to Chappe Private Messaging for EOSIO\n\n"

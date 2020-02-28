package cmd

import (
	"fmt"

	"github.com/morikuni/aec"
	"github.com/spf13/cobra"
)

var (
	Version   string
	GitCommit string
)

func PrintChappeASCIIArt() {
	chappeLogo := aec.BlueF.Apply(chappeFigletStr)
	fmt.Print(chappeLogo)
}

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
const chappeFigletStr = 'chappe'
//  	 _
//  ___| |__   __ _ _ __  _ __   ___
// / __| '_ \ / _` | '_ \| '_ \ / _ \
// | (__| | | | (_| | |_) | |_) |  __/
// \___|_| |_|\__,_| .__/| .__/ \___|
// 				|_|   |_|

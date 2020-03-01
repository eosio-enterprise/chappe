package cmd

import (
	"bytes"
	"fmt"
	"os"

	shell "github.com/ipfs/go-ipfs-api"

	"github.com/spf13/cobra"
)

// MakeGet ...
func MakeGet() *cobra.Command {
	var command = &cobra.Command{
		Use:   "get",
		Short: "Get chappe artifacts",
		Long:  `Retrieve and optionally decrypt a specific item from the data store`,
		Example: `  chappe get
  chappe get --cid Qmah6HuF5kw1HF9yNfTxxoWBjCRYddestuMYn7RFF4RHnS`,
		SilenceUsage: false,
	}

	command.Flags().StringP("cid", "c", "", "CID hash uniquely identifying object in data store")

	command.RunE = func(command *cobra.Command, args []string) error {

		cid, _ := command.Flags().GetString("cid")
		if len(cid) == 0 {
			return fmt.Errorf("--cid required")
		}

		sh := shell.NewShell("localhost:5001") // TODO: move to configuration
		reader, err := sh.Cat(cid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not add encrypted data to IPFS: %s", err)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		catRand := buf.String()

		fmt.Println("Object Contents:")
		fmt.Println(catRand)

		return nil
	}

	return command
}

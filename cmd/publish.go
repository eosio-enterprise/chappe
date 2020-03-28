package cmd

import (
	"encoding/json"
	"log"
	"time"

	"fmt"

	"github.com/eosio-enterprise/chappe/internal/encryption"
	"github.com/eosio-enterprise/chappe/pkg"
	"github.com/fatih/color"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// MakePublish ...
func MakePublish() *cobra.Command {
	var command = &cobra.Command{
		Use:   "publish",
		Short: "Publish a private message to a channel",
		Long: `Create and share private messages with other nodes on a chappe network.
					Currently, only a randomly generated JSON object is sent as the
					message body. Other options will be added shortly.`,
		Example:      `chappe publish --channel-name MyChannel --memo "This memo is not encrypted"`,
		SilenceUsage: false,
	}

	var channelName, readableMemo string
	var encryptFlag bool
	command.Flags().StringVarP(&channelName, "channel-name", "n", "", "channel name")
	command.Flags().StringVarP(&readableMemo, "memo", "m", "", "Human readable memo to attach to payload (never encrypted)")
	command.Flags().BoolVarP(&encryptFlag, "encrypt", "", true, "Boolean flag whether to encrypt payload - defaults to true")

	command.RunE = func(command *cobra.Command, args []string) error {

		if encryptFlag && len(channelName) == 0 {
			return fmt.Errorf("--channel-name is required when encrypting")
		}

		for {
			privatePayload := pkg.GetLedgerPayload()
			msg := pkg.NewMessage()
			msg.Payload["BlockchainMemo"] = []byte(readableMemo)

			payload, _ := json.MarshalIndent(privatePayload, "", "  ")
			color.Green("Publishing:")
			color.Green(string(payload))

			if encryptFlag {
				aesKey := encryption.NewAesEncryptionKey()
				aesEncryptedData, err := encryption.AesEncrypt(payload, aesKey)
				if err != nil {
					log.Panicf("Error with AES encryption: %s", err)
				}

				msg.Payload["EncryptedPayload"] = aesEncryptedData
				encryptedAesKey, err := encryption.RsaEncrypt(channelName, aesKey[:])
				if err != nil {
					log.Panicf("Error with RSA encryption: %s", err)
				}
				msg.Payload["EncryptedAESKey"] = encryptedAesKey
			}

			trxID, err := pkg.Publish(msg)
			if err != nil {
				log.Println("Error submitting transaction to EOSIO: ", err)
			}

			log.Println("Published to channel: ", channelName, "; TrxId: "+trxID)

			time.Sleep(viper.GetDuration("PublishInterval"))
		}
	}
	return command
}

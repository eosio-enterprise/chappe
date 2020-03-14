package cmd

import (
	"encoding/json"
	"log"
	"time"

	"fmt"

	"github.com/eosio-enterprise/chappe/internal/encryption"
	"github.com/eosio-enterprise/chappe/pkg"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var c = cache.New(time.Hour)

// MakePublish ...
func MakePublish() *cobra.Command {
	var command = &cobra.Command{
		Use:   "publish",
		Short: "Publish a private message to a channel",
		Long:  `Create and share private messages with other nodes on a chappe network.`,
		Example: `  chappe publish
  chappe publish --channel-name MyChannel --readable-memo "This memo is not encrypted"
  chappe publish --encrypt false`,

		SilenceUsage: false,
	}

	var channelName, readableMemo string
	var encryptFlag bool

	command.Flags().StringVarP(&channelName, "channel-name", "n", "", "channel name")
	command.Flags().StringVarP(&readableMemo, "readable-memo", "m", "", "Human readable memo to attach to payload (never encrypted)")
	command.Flags().BoolVarP(&encryptFlag, "encrypt", "", true, "Boolean flag whether to encrypt payload - defaults to true")

	command.RunE = func(command *cobra.Command, args []string) error {

		if encryptFlag && len(channelName) == 0 {
			return fmt.Errorf("--channel-name is required when encrypting")
		}

		for {
			privatePayload := pkg.GetFakePrivatePayload()
			msg := pkg.NewMessage()
			msg.Payload["BlockchainMemo"] = []byte(readableMemo)

			payload, _ := json.MarshalIndent(privatePayload, "", "  ")
			log.Println("Publishing: \n", string(payload))

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

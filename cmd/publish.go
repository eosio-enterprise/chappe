package cmd

import (
	"encoding/hex"
	"encoding/json"
	"strings"

	"fmt"
	"os"

	"github.com/bxcodec/faker"
	"github.com/eoscanada/eos-go"
	"github.com/eosio-enterprise/chappe/cmd/message"
	"github.com/eosio-enterprise/chappe/internal/encryption"
	shell "github.com/ipfs/go-ipfs-api"

	// "github.com/polydawn/refmt/json"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// PrivatePayload ...
type PrivatePayload struct {
	RecordID         string `faker:"uuid_hyphenated"`
	FirstName        string `faker:"first_name"`
	LastName         string `faker:"last_name"`
	DOB              string `faker:"date"`
	CreditCardNumber string `faker:"cc_number"`
	CreditCardType   string `faker:"cc_type"`
	Email            string `faker:"email"`
	TimeZone         string `faker:"timezone"`
	AmountDue        string `faker:"amount_with_currency"`
	PhoneNumber      string `faker:"phone_number"`
	SafeWord         string `faker:"word"`
	LastScan         string `faker:"timestamp"`
}

// PersistedObject ...
type PersistedObject struct {
	EncryptedPayload   []byte
	EncryptedAESKey    []byte
	UnencryptedPayload string
}

func publish(eosioEndpoint, channelName string, payload PersistedObject) (string, error) {
	api := eos.New(eosioEndpoint)

	keyBag := &eos.KeyBag{}
	err := keyBag.ImportPrivateKey(viper.GetString("Eosio.PublishPrivateKey"))
	if err != nil {
		panic(fmt.Errorf("import private key: %s", err))
	}
	api.SetSigner(keyBag)

	sh := shell.NewShell(viper.GetString("IPFS.Endpoint")) // TODO: move to configuration
	jsonPayloadNode, err := json.Marshal(payload)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not marshall:  %s", err)
	}

	hash, err := sh.Add(strings.NewReader(string(jsonPayloadNode)))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not add data to IPFS: %s", err)
	}
	fmt.Println("IPFS Hash: ", hash)

	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(api); err != nil {
		panic(fmt.Errorf("filling tx opts: %s", err))
	}

	tx := eos.NewTransaction([]*eos.Action{message.NewPub(hash, readableMemo)}, txOpts)
	_, packedTx, err := api.SignTransaction(tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		panic(fmt.Errorf("sign transaction: %s", err))
	}

	response, err := api.PushTransaction(packedTx)
	if err != nil {
		panic(fmt.Errorf("push transaction: %s", err))
	}
	return hex.EncodeToString(response.Processed.ID), nil
}

var channelName string
var eosioEndpoint string
var readableMemo string
var encryptFlag bool

// MakePublish ...
func MakePublish() *cobra.Command {
	var command = &cobra.Command{
		Use:   "publish",
		Short: "Publish a private message to a channel",
		Long:  `Create and share private messages with other nodes on a chappe network.`,
		Example: `  chappe publish
  chappe publish --channel-name MyChannel --eosio-endpoint https://eos.greymass.com
  chappe publish --encrypt false --eosio-endpoint https://jungle2.cryptolions.io`,

		SilenceUsage: false,
	}

	command.Flags().StringVarP(&channelName, "channel-name", "n", "default", "channel name")
	command.Flags().StringVarP(&eosioEndpoint, "eosio-endpoint", "e", "https://jungle2.cryptolions.io", "EOSIO JSONRPC endpoint")
	command.Flags().StringVarP(&readableMemo, "readable-memo", "m", "", "Human readable memo to attach to payload (never encrypted)")
	command.Flags().BoolVarP(&encryptFlag, "encrypt", "", true, "Boolean flag whether to encrypt payload - defaults to true")

	command.RunE = func(command *cobra.Command, args []string) error {

		if encryptFlag && len(channelName) == 0 {
			return fmt.Errorf("--channel-name is required when encrypting")
		}

		privatePayload := PrivatePayload{}
		err := faker.FakeData(&privatePayload)
		if err != nil {
			fmt.Println(err)
		}

		persistedObject := PersistedObject{}
		persistedObject.UnencryptedPayload = readableMemo

		payload, _ := json.MarshalIndent(privatePayload, "", "  ")
		fmt.Println("Broadcasting Message:")
		fmt.Println(string(payload))

		if encryptFlag {
			aesKey := encryption.NewAesEncryptionKey()
			aesEncryptedData, err := encryption.AesEncrypt(payload, aesKey)
			if err != nil {
				panic(fmt.Errorf("Error trying to encrypt %s", err))
			}

			persistedObject.EncryptedPayload = aesEncryptedData

			encryptedAesKey, err := encryption.RsaEncrypt(channelName, aesKey[:])
			if err != nil {
				panic(fmt.Errorf("Error encrypting the AES key %s", err))
			}
			persistedObject.EncryptedAESKey = encryptedAesKey
		}

		trxID, err := publish(eosioEndpoint, channelName, persistedObject)
		if err != nil {
			fmt.Println("Error submitting transaction to EOSIO: ", err)
		}

		fmt.Println(
			`=======================================================================
Published message to channel ` + channelName + `
Transaction ID: ` + trxID + `
=======================================================================`)

		return nil
	}

	return command
}

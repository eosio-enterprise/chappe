package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"unsafe"

	"github.com/bxcodec/faker"
	"github.com/eoscanada/eos-go"
	"github.com/eosio-enterprise/chappe/cmd/message"
	"github.com/spf13/cobra"
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
	EncryptedPayload string
	EncryptedAESKey  string
	readableMemo     string
}

// ParseRsaPublicKeyFromPem ...
func ParseRsaPublicKeyFromPem(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}

// ParseRsaPrivateKeyFromPem ...
func ParseRsaPrivateKeyFromPem(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func loadPublicKey(keyname string) *rsa.PublicKey {

	publicKeyPemStr, err := ioutil.ReadFile(keyname + ".pub")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from reading key file: %s\n", err)
		return nil
	}

	publicKey, err := ParseRsaPublicKeyFromPem(publicKeyPemStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from reading key file: %s\n", err)
		return nil
	}

	return publicKey
}

func rsaEncrypt(channelName, payload string) ([]byte, error) {
	publicKey := loadPublicKey(channelName)

	secretMessage, _ := json.Marshal(payload)
	label := []byte("chappe") // TODO: migrate to something else?

	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, secretMessage, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
		return nil, err
	}
	return encryptedData, nil
}

func readPrivateKey() string {
	// Right now, the key is read from an environment variable, it's an example after all.
	// In a real-world scenario, would you probably integrate with a real wallet or something similar
	envName := "EOS_GO_PRIVATE_KEY"
	privateKey := os.Getenv(envName)
	if privateKey == "" {
		panic(fmt.Errorf("private key environment variable %q must be set", envName))
	}

	return privateKey
}

func publish(eosioEndpoint, channelName, readableMemo string, payload []byte) (string, error) {
	api := eos.New(eosioEndpoint)

	keyBag := &eos.KeyBag{}
	err := keyBag.ImportPrivateKey(readPrivateKey())
	if err != nil {
		panic(fmt.Errorf("import private key: %s", err))
	}
	api.SetSigner(keyBag)

	// sh := shell.NewShell("localhost:5001") // TODO: move to configuration
	cid := "cid"
	// err := sh.Add(strings.NewReader(string(secretMessage))) // TODO: does this really need to be using "strings"
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "Could not add encrypted data to IPFS: %s", err)
	// }
	fmt.Println("IPFS Hash: ", cid)

	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(api); err != nil {
		panic(fmt.Errorf("filling tx opts: %s", err))
	}

	tx := eos.NewTransaction([]*eos.Action{message.NewPub(cid, readableMemo)}, txOpts)
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

func bytesToString(b [32]byte) string {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := reflect.StringHeader{bh.Data, bh.Len}
	return *(*string)(unsafe.Pointer(&sh))
}

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
		// eosioEndpoint, _ := command.Flags().GetString("eosio-endpoint")
		// channelName, _ := command.Flags().GetString("channel-name")
		// readableMemo, _ := command.Flags().GetString("readable-memo")
		// encryptFlag, _ := command.Flags().GetString("encyrpt")

		// fmt.Println("ChannelName	: ", channelName)
		// fmt.Println("EosioEndpoint	: ", eosioEndpoint)
		// fmt.Println("ReadableMemo	: ", readableMemo)
		// fmt.Println("Encrypt		: ", encryptFlag)

		if encryptFlag && len(channelName) == 0 {
			return fmt.Errorf("--channel-name is required when encrypting")
		}

		payload := PrivatePayload{}
		err := faker.FakeData(&payload)
		if err != nil {
			fmt.Println(err)
		}

		jsonPayload, _ := json.MarshalIndent(payload, "", "  ")
		fmt.Println("Unencrypted Data:")
		fmt.Println(string(jsonPayload))

		if encryptFlag {
			aesKey := NewEncryptionKey()
			fmt.Println("aesKey			: ", aesKey)
			fmt.Println("aesKey string	: ", string(aesKey[:]))

			// aesEncryptedData, err := SymmetricEncrypt(aesKey, jsonPayload)
			aesEncryptedData, err := Encrypt(jsonPayload, aesKey)
			if err != nil {
				panic(fmt.Errorf("Error trying to encrypt %s", err))
			}
			fmt.Println("AES Encrypted Data:")
			fmt.Println(aesEncryptedData)

			// aesUnencryptedData, err := SymmetricDecrypt(aesKey, aesEncryptedData)
			aesUnencryptedData, err := Decrypt(aesEncryptedData, aesKey)
			if err != nil {
				panic(fmt.Errorf("Error trying to dencrypt %s", err))
			}
			fmt.Println("AES Decrypted Data:")
			fmt.Println(string(aesUnencryptedData))

		}

		//     	trxID, err := publish(eosioEndpoint, channelName, readableMemo, payload)
		// 		if err != nil {
		// 			fmt.Println("Error submitting transaction to EOSIO: ", err)
		// 		}

		// 		fmt.Println(
		// 			`=======================================================================
		// Published message to channel ` + channelName + `
		// Transaction ID: ` + trxID + `
		// =======================================================================`)

		return nil
	}

	return command
}

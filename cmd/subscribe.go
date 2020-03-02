package cmd

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dfuse-io/eosws-go"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

// MakeSubscribe ...
func MakeSubscribe() *cobra.Command {
	var command = &cobra.Command{
		Use:   "subscribe",
		Short: "Subscribe to a channel",
		Example: `  chappe subscribe
	chappe subscribe --channel-name MyChannel --eosio-endpoint https://eos.greymass.com`,
		SilenceUsage: false,
	}

	command.Flags().StringVarP(&channelName, "channel-name", "n", "default", "channel name")
	command.Flags().StringVarP(&eosioEndpoint, "eosio-endpoint", "e", "https://jungle2.cryptolions.io", "EOSIO JSONRPC endpoint")

	command.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("Ctrl-C to exit")
		fmt.Println("Contribute: https://github.com/eosio-enterprise/chappe/blob/master/CONTRIBUTING.md")

		go func() {

			client := getClient()
			err := client.Send(getActionTraces())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Client cannot send request: %s", err)
			}

			for {
				msg, err := client.Read()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Cannot read from client: %s", err)
				}

				switch m := msg.(type) {
				case *eosws.ActionTrace:

					ipfsHash := gjson.Get(string(m.Data.Trace), "act.data.ipfs_hash")
					memo := gjson.Get(string(m.Data.Trace), "act.data.memo")

					fmt.Println("IPFS Hash :", ipfsHash)
					fmt.Println("Memo		:", memo)

					sh := shell.NewShell("localhost:5001")
					reader, err := sh.Cat(ipfsHash.String())
					if err != nil {
						fmt.Fprintf(os.Stderr, "Could not not find IPFS hash %s", err)
					}

					buf := new(bytes.Buffer)
					buf.ReadFrom(reader)
					retrievedContents := buf.String()

					var persistedObject PersistedObject
					if err := json.Unmarshal([]byte(retrievedContents), &persistedObject); err != nil {
						panic(err)
					}

					aesKey, err := rsaDecrypt(channelName, persistedObject.EncryptedAESKey)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Cannot decrypt the AES key: %s", err)
					}

					plaintext, err := aesDecrypt(persistedObject.EncryptedPayload, aesKey)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error from decryption: %s\n", err)
					}

					var receivedPrivatePayload PrivatePayload
					err = json.Unmarshal(plaintext, &receivedPrivatePayload)

					indentedPayload, err := json.MarshalIndent(receivedPrivatePayload, "", "  ")
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error from marshall indent: %s\n", err)
					}

					fmt.Printf("Received message: %s\n", string(indentedPayload))

				case *eosws.Progress:
					fmt.Println("Received Progress Message at Block:", m.Data.BlockNum)
				case *eosws.Listening:
					fmt.Println("Received Listening Message ...")
				default:
					fmt.Println("Unsupported message", m)
				}
			}
		}()

		sigs := make(chan os.Signal, 1)
		<-sigs
	}
	return command
}

var dfuseEndpoint = "wss://jungle.eos.dfuse.io/v1/stream"
var origin = "https://origin.example.io"

func parseRsaPrivateKeyFromPem(privPEM []byte) (*rsa.PrivateKey, error) {
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

func load(keyname string) *rsa.PrivateKey {

	privateKeyPemStr, err := ioutil.ReadFile("" + keyname + ".pem")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from reading key file: %s\n", err)
		return nil
	}

	priv, err := parseRsaPrivateKeyFromPem(privateKeyPemStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from reading key file: %s\n", err)
		return nil
	}
	return priv
}

func rsaDecrypt(channelName string, payload []byte) ([]byte, error) {
	privateKey := load(channelName)
	label := []byte("chappe") // TODO: migrate to something else?

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, payload, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot decrypt the AES key: %s\n", err)
		return nil, err
	}

	return plaintext, nil
}

// Decrypt decrypts data using 256-bit AES-GCM.  This both hides the content of
// the data and provides a check that it hasn't been altered. Expects input
// form nonce|ciphertext|tag where '|' indicates concatenation.
func aesDecrypt(ciphertext []byte, key []byte) (plaintext []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("malformed ciphertext")
	}

	return gcm.Open(nil,
		ciphertext[:gcm.NonceSize()],
		ciphertext[gcm.NonceSize():],
		nil,
	)
}

func getToken(apiKey string) (token string, expiration time.Time, err error) {
	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{"api_key":"%s"}`, apiKey)))
	resp, err := http.Post("https://auth.dfuse.io/v1/auth/issue", "application/json", reqBody)
	if err != nil {
		err = fmt.Errorf("unable to obtain token: %s", err)
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("unable to obtain token, status not 200, got %d: %s", resp.StatusCode, reqBody.String())
		return
	}

	if body, err := ioutil.ReadAll(resp.Body); err == nil {
		token = gjson.GetBytes(body, "token").String()
		expiration = time.Unix(gjson.GetBytes(body, "expires_at").Int(), 0)
	}
	return
}

func getClient() *eosws.Client {
	apiKey := os.Getenv("DFUSE_API_KEY")
	if apiKey == "" {
		log.Fatalf("Set your API key to environment variable DFUSE_API_KEY")
	}

	jwt, _, err := eosws.Auth(apiKey)
	if err != nil {
		log.Fatalf("cannot get auth token: %s", err.Error())
	}

	dfuseEndpoint := "wss://jungle.eos.dfuse.io/v1/stream"
	origin := "github.com/eosio-enterprise/chappe"
	client, err := eosws.New(dfuseEndpoint, jwt, origin)
	if err != nil {
		log.Fatalf("cannot connect to dfuse endpoint: %s", err.Error())
	}
	return client
}

func getActionTraces() *eosws.GetActionTraces {
	ga := &eosws.GetActionTraces{}
	ga.ReqID = "chappe GetActions"
	ga.StartBlock = -300
	ga.Listen = true
	ga.WithProgress = 3
	ga.IrreversibleOnly = false
	ga.Data.Accounts = "messengerbus"
	ga.Data.ActionNames = "pub"
	fmt.Printf("Connecting to network...  %s::%s\n", ga.Data.Accounts, ga.Data.ActionNames)
	ga.Data.WithInlineTraces = true
	return ga
}

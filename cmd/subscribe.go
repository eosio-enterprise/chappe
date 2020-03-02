package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dfuse-io/eosws-go"
	"github.com/eosio-enterprise/chappe/internal/encryption"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	command.Flags().StringVarP(&eosioEndpoint, "eosio-endpoint", "e", viper.GetString("Eosio.Endpoint"), "EOSIO JSONRPC endpoint")

	command.Run = func(cmd *cobra.Command, args []string) {

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

					sh := shell.NewShell(viper.GetString("IPFS.Endpoint"))
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

					aesKey, err := encryption.RsaDecrypt(channelName, persistedObject.EncryptedAESKey)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Cannot decrypt the AES key: %s", err)
					}

					plaintext, err := encryption.AesDecrypt(persistedObject.EncryptedPayload, &aesKey)
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
	apiKey := viper.GetString("Dfuse.ApiKey")
	if apiKey == "" {
		log.Fatalf("Missing Dfuse.ApiKey in config")
	}

	jwt, _, err := eosws.Auth(apiKey)
	if err != nil {
		log.Fatalf("cannot get auth token: %s", err.Error())
	}

	// var dfuseEndpoint = viper.GetString("Dfuse.WSEndpoint")
	var origin = viper.GetString("Dfuse.Origin")

	client, err := eosws.New(viper.GetString("Dfuse.WSEndpoint"), jwt, origin)
	if err != nil {
		log.Fatalf("cannot connect to dfuse endpoint: %s", err.Error())
	}
	return client
}

func getActionTraces() *eosws.GetActionTraces {
	ga := &eosws.GetActionTraces{}
	ga.ReqID = "chappe"
	ga.StartBlock = -300
	ga.Listen = true
	ga.WithProgress = 3
	ga.IrreversibleOnly = false
	ga.Data.Accounts = viper.GetString("Eosio.PublishAccount")
	ga.Data.ActionNames = "pub"
	fmt.Printf("Connecting...  %s::%s\n", ga.Data.Accounts, ga.Data.ActionNames)
	ga.Data.WithInlineTraces = true
	return ga
}

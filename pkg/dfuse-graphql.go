package pkg

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	pb "github.com/eosio-enterprise/chappe/internal/pb"
	"github.com/spf13/viper"

	"github.com/tidwall/gjson"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func getToken(apiKey string) (token string, expiration time.Time, err error) {
	reqBody := bytes.NewBuffer([]byte(fmt.Sprintf(`{"api_key":"%s"}`, apiKey)))
	resp, err := http.Post(viper.GetString("Dfuse.AuthEndpoint"), "application/json", reqBody)
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

func createClient(endpoint string) pb.GraphQLClient {
	dfuseAPIKey := viper.GetString("Dfuse.ApiKey")
	if dfuseAPIKey == "" {
		log.Fatal("Dfuse.ApiKey is required in configuration")
	}

	token, _, err := getToken(dfuseAPIKey)
	if err != nil {
		log.Fatalf("Cannot retrieve dfuse token: %s", err.Error())
	}

	credential := oauth.NewOauthAccess(&oauth2.Token{AccessToken: token, TokenType: "Bearer"})
	transportCreds := credentials.NewClientTLSFromCert(nil, "")
	conn, err := grpc.Dial(endpoint,
		grpc.WithPerRPCCredentials(credential),
		grpc.WithTransportCredentials(transportCreds),
	)
	if err != nil {
		log.Fatalf("Cannot dial to endpoint: %s", err.Error())
	}

	return pb.NewGraphQLClient(conn)
}

// TODO: need to replace messengerbus with viper.GetString("Eosio.PublishAccount")
const operationEOS = `subscription {
	searchTransactionsForward(query:"receiver:messengerbus action:pubmap") {
	  undo cursor
	  trace { id matchingActions { json } }
	}
  }`

type eosioDocument struct {
	SearchTransactionsForward struct {
		Cursor string
		Undo   bool
		Trace  struct {
			ID              string
			MatchingActions []struct {
				JSON map[string]interface{}
			}
		}
	}
}

// StreamMessages ...
func StreamMessages(ctx context.Context, channelName, ledgerFile string, sendReceipts bool) {
	/* The client can be re-used for all requests, cache it at the appropriate level */
	client := createClient(viper.GetString("Dfuse.GraphQLEndpoint"))
	executor, err := client.Execute(ctx, &pb.Request{Query: operationEOS})
	if err != nil {
		log.Fatalf("Cannot execute dfuse graphql query: %s", err.Error())
	}

	for {
		resp, err := executor.Recv()
		if err != nil {
			log.Fatalf("Cannot recv on dfuse graphql: %s", err.Error())
		}

		if len(resp.Errors) > 0 {
			for _, err := range resp.Errors {
				log.Printf("Request failed: %s\n", err)
			}

			/* We continue here, but you could take another decision here, like exiting the process */
			continue
		}

		document := &eosioDocument{}
		err = json.Unmarshal([]byte(resp.Data), document)
		if err != nil {
			log.Fatalf("Cannot unmarshal dfuse graphql document: %s", err.Error())
		}

		result := document.SearchTransactionsForward
		if result.Undo {
			log.Println("EOSIO transaction has been reverted, halting process. Skipping")
			continue
		} else {
			for _, action := range result.Trace.MatchingActions {
				data := action.JSON["payload"].([]interface{})
				payload := make(map[string]string)
				for i := 0; i < len(data); i++ {
					imap := data[i].(map[string]interface{})
					payload[imap["key"].(string)] = imap["value"].(string)
				}

				if payload["message_type"] == "receipt" {
					log.Println("Received receipt, ignoring.")
				} else {
					message, err := receiveGQL(channelName, ledgerFile, payload)
					if err == nil && sendReceipts {
						SendReceipt(channelName, message)
					}
				}
			}
		}
	}
}

func receiveGQL(channelName, ledgerFile string, payload map[string]string) (Message, error) {
	var msg Message
	cid, cidExists := payload["cid"]
	if cidExists {
		log.Println("Received notification of new message: ", payload["cid"], "; memo: ", payload["memo"])
		msg, err := Load(channelName, ledgerFile, cid)
		if err != nil {
			log.Println("Error loading message: ", err)
			return msg, err
		}
		return msg, nil
	}
	log.Println("Message does not contain a CID, discarding.")
	return msg, nil
}

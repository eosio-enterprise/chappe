package pkg

import (
	"fmt"
	"log"

	"github.com/dfuse-io/eosws-go"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// StreamWS ...
func StreamWS(channelName string, sendReceipts bool) {
	client := getClient()
	err := client.Send(getActionTraces())
	if err != nil {
		log.Fatalf("Failed to send request to dfuse: %s", err)
	}

	for {
		msg, err := client.Read()
		if err != nil {
			log.Fatalf("Cannot read from dfuse client: %s", err)
		}

		switch m := msg.(type) {
		case *eosws.ActionTrace:
			message, err := receiveWS(channelName, m)
			if err == nil && sendReceipts {
				SendReceipt(channelName, message)
			}
		case *eosws.Progress:
			fmt.Print(".") // poor man's progress bar, using print not log
		case *eosws.Listening:
			log.Println("Received Listening Message ...")
		default:
			log.Println("Received Unsupported Message", m)
		}
	}
}

func receiveWS(channelName string, dfuseMessage *eosws.ActionTrace) (Message, error) {
	ipfsHash := gjson.Get(string(dfuseMessage.Data.Trace), "act.data.ipfs_hash")
	memo := gjson.Get(string(dfuseMessage.Data.Trace), "act.data.memo")
	fmt.Println()
	log.Println("Received notification of new message: ", ipfsHash, "; memo: ", memo)

	msg, err := Load(channelName, ipfsHash.String())
	if err != nil {
		log.Println("Error loading message: ", err)
		return msg, err
	}

	return msg, nil
}

// GetClient ...
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

// GetActionTraces ...
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

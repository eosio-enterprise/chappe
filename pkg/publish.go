package pkg

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"

	"github.com/eoscanada/eos-go"
	"github.com/spf13/viper"
)

// func newPub(ipfsHash, memo string) *eos.Action {
// 	return &eos.Action{
// 		Account: eos.AN(viper.GetString("Eosio.PublishAccount")),
// 		Name:    eos.ActN("pub"),
// 		Authorization: []eos.PermissionLevel{
// 			{Actor: eos.AN(viper.GetString("Eosio.PublishAccount")), Permission: eos.PN("active")},
// 		},
// 		ActionData: eos.NewActionData(PubActionPayload{
// 			IpfsHash: ipfsHash,
// 			Memo:     memo,
// 		}),
// 	}
// }

// PubMapActionPayload ...
type PubMapActionPayload struct {
	PayloadMap map[string]string `json:"payload"`
}

func newPubMap(payload map[string]string) *eos.Action {
	return &eos.Action{
		Account: eos.AN(viper.GetString("Eosio.PublishAccount")),
		Name:    eos.ActN("pubmap"),
		Authorization: []eos.PermissionLevel{
			{Actor: eos.AN(viper.GetString("Eosio.PublishAccount")), Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(PubMapActionPayload{
			PayloadMap: payload,
		}),
	}
}

// // PubActionPayload ...
// type PubActionPayload struct {
// 	IpfsHash string `json:"ipfs_hash"`
// 	Memo     string `json:"memo"`
// }

// PublishMapToBlockchain ...
func PublishMapToBlockchain(payload map[string]string) (string, error) {
	ctx := context.Background()
	api := eos.New(viper.GetString("Eosio.Endpoint"))

	keyBag := &eos.KeyBag{}
	err := keyBag.ImportPrivateKey(ctx, viper.GetString("Eosio.PublishPrivateKey"))
	if err != nil {
		log.Panicf("import private key: %s", err)
	}
	api.SetSigner(keyBag)

	txOpts := &eos.TxOptions{}
	if err := txOpts.FillFromChain(ctx, api); err != nil {
		log.Printf("Error filling tx opts: %s", err)
		return "error", err
	}

	tx := eos.NewTransaction([]*eos.Action{newPubMap(payload)}, txOpts)
	_, packedTx, err := api.SignTransaction(ctx, tx, txOpts.ChainID, eos.CompressionNone)
	if err != nil {
		log.Printf("Error signing transaction: %s", err)
		return "error", err
	}

	response, err := api.PushTransaction(ctx, packedTx)
	if err != nil {
		log.Printf("Error pushing transaction: %s", err)
		return "error", err
	}
	return hex.EncodeToString(response.Processed.ID), nil
}

// // AddToEosio ...
// func AddToEosio(cid string, readableMemo string) (string, error) {
// 	ctx := context.Background()
// 	api := eos.New(viper.GetString("Eosio.Endpoint"))

// 	keyBag := &eos.KeyBag{}
// 	err := keyBag.ImportPrivateKey(ctx, viper.GetString("Eosio.PublishPrivateKey"))
// 	if err != nil {
// 		log.Panicf("import private key: %s", err)
// 	}
// 	api.SetSigner(keyBag)

// 	txOpts := &eos.TxOptions{}
// 	if err := txOpts.FillFromChain(ctx, api); err != nil {
// 		log.Printf("Error filling tx opts: %s", err)
// 		return "error", err
// 	}

// 	tx := eos.NewTransaction([]*eos.Action{newPub(cid, readableMemo)}, txOpts)
// 	_, packedTx, err := api.SignTransaction(ctx, tx, txOpts.ChainID, eos.CompressionNone)
// 	if err != nil {
// 		log.Printf("Error signing transaction: %s", err)
// 		return "error", err
// 	}

// 	response, err := api.PushTransaction(ctx, packedTx)
// 	if err != nil {
// 		log.Printf("Error pushing transaction: %s", err)
// 		return "error", err
// 	}
// 	return hex.EncodeToString(response.Processed.ID), nil
// }

func addToIpfs(payload Message) string {
	sh := shell.NewShell(viper.GetString("IPFS.Endpoint"))
	jsonPayloadNode, err := json.Marshal(payload)
	if err != nil {
		log.Panicf("Could not marshall:  %s", err)
	}

	cid, err := sh.Add(strings.NewReader(string(jsonPayloadNode)))
	if err != nil {
		log.Printf("Could not add data to IPFS: %s", err)
	}
	log.Println("Saved message to IPFS; CID: ", cid)
	return cid
}

// Publish ...
func Publish(payload Message) (string, error) {
	blockchainMemo, memoExists := payload.Payload["BlockchainMemo"]
	if !memoExists {
		blockchainMemo = []byte("")
	}

	eosioPayload := make(map[string]string)
	eosioPayload["cid"] = addToIpfs(payload)
	eosioPayload["memo"] = string(blockchainMemo)
	return PublishMapToBlockchain(eosioPayload)
}

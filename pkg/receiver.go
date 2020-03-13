package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/dfuse-io/eosws-go"
	"github.com/eosio-enterprise/chappe/internal/encryption"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

// Receive ...
func Receive(channelName string, dfuseMessage *eosws.ActionTrace) error {
	ipfsHash := gjson.Get(string(dfuseMessage.Data.Trace), "act.data.ipfs_hash")
	memo := gjson.Get(string(dfuseMessage.Data.Trace), "act.data.memo")
	fmt.Println()
	log.Println("Received notification of new message: ", ipfsHash, "; memo: ", memo)

	sh := shell.NewShell(viper.GetString("IPFS.Endpoint"))
	reader, err := sh.Cat(ipfsHash.String())
	if err != nil {
		log.Println("Could not not find IPFS hash: ", ipfsHash, "; Error: ", err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	retrievedContents := buf.String()

	var msg Message
	if err := json.Unmarshal([]byte(retrievedContents), &msg); err != nil {
		log.Println("Cannot unmarshal object retrieved from IPFS: ", err)
	}

	aesKey, err := encryption.RsaDecrypt(channelName, msg.EncryptedAESKey)
	if err != nil {
		trxID := string(dfuseMessage.Data.TransactionID)
		log.Println("Cannot read contents of message, discarding. TrxId: ", trxID)
	} else {
		plaintext, err := encryption.AesDecrypt(msg.EncryptedPayload, &aesKey)
		if err != nil {
			log.Printf("Error from AES decryption: %s\n", err)
		}

		log.Println("Decrypted message from channelName: \n", string(plaintext))

		// We could also unmarshal the text (which is JSON) into an object
		// var receivedPrivatePayload FakePrivatePayload
		// err = json.Unmarshal(plaintext, &receivedPrivatePayload)

		// indentedPayload, err := json.MarshalIndent(receivedPrivatePayload, "", "  ")
		// if err != nil {
		// 	log.Printf("Error from marshall indent: %s\n", err)
		// }
	}
	return nil
}

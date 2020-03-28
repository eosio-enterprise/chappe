package pkg

import (
	"encoding/hex"
	"log"

	"github.com/eosio-enterprise/chappe/internal/encryption"
)

// SendReceipt ...
func SendReceipt(channelName string, msg Message) error {

	receiptSignature, err := encryption.CalcReceipt(msg.Bytes())
	if err != nil {
		// TODO: improve handling if the receipt cannot be generated and sent, perhaps should be fatal
		log.Println("Cannot send receipt: ", err)
	}
	receiptStr := hex.EncodeToString(receiptSignature)
	//log.Println("Sending receipt: ", receiptStr)

	receiptMap := make(map[string]string)
	receiptMap["receipt"] = receiptStr
	receiptMap["message_type"] = "receipt"

	// currently, receipts will reveal other parties on the channel
	// TODO: mask metadata on receipts
	trxID, _ := PublishMapToBlockchain(receiptMap)

	log.Println("Sent receipt, transaction ID: ", trxID)
	return nil
}

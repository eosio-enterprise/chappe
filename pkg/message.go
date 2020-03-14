package pkg

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/bxcodec/faker"
	"github.com/eosio-enterprise/chappe/internal/encryption"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/spf13/viper"
)

// Message ...
type Message struct {
	EncryptedPayload []byte
	EncryptedAESKey  []byte
	ReadableMemos    []string
	Payload          map[string][]byte
}

// NewMessage - creates a new Message object, provisions the Payload map
func NewMessage() Message {
	m := Message{}
	m.Payload = make(map[string][]byte)
	return m
}

// Load ...
func Load(channelName, ipfsHash string) (Message, error) {
	var msg Message
	sh := shell.NewShell(viper.GetString("IPFS.Endpoint"))
	reader, err := sh.Cat(ipfsHash)
	if err != nil {
		log.Println("Could not not find IPFS hash: ", ipfsHash, "; Error: ", err)
		return msg, err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(reader)
	// retrievedContents := buf.String()

	if err := json.Unmarshal(buf.Bytes(), &msg); err != nil {
		log.Println("Cannot unmarshal object retrieved from IPFS: ", err)
		return msg, err
	}

	aesKey, err := encryption.RsaDecrypt(channelName, msg.Payload["EncryptedAESKey"])
	if err != nil {
		log.Println("Cannot read contents of message, discarding. CID: ", ipfsHash)
		return msg, err
	}

	plaintext, err := encryption.AesDecrypt(msg.Payload["EncryptedPayload"], &aesKey)
	if err != nil {
		log.Printf("Error from AES decryption: %s\n", err)
		return msg, err
	}

	log.Println("Decrypted message from Channel: ", channelName, "\n", string(plaintext))

	// We could also unmarshal the text (which is JSON) into an object
	// var receivedPrivatePayload FakePrivatePayload
	// err = json.Unmarshal(plaintext, &receivedPrivatePayload)

	// indentedPayload, err := json.MarshalIndent(receivedPrivatePayload, "", "  ")
	// if err != nil {
	// 	log.Printf("Error from marshall indent: %s\n", err)
	// }

	return msg, nil
}

// FakePrivatePayload ...
type FakePrivatePayload struct {
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

// GetFakePrivatePayload ...
func GetFakePrivatePayload() FakePrivatePayload {
	var privatePayload FakePrivatePayload
	err := faker.FakeData(&privatePayload)
	if err != nil {
		log.Println("Cannot generate fake data: ", err)
	}
	return privatePayload
}

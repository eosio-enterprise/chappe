package pkg

import (
	"log"

	"github.com/bxcodec/faker"
)

// Message ...
type Message struct {
	EncryptedPayload   []byte
	EncryptedAESKey    []byte
	UnencryptedPayload string
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

package pkg

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/bxcodec/faker/v3"
	"github.com/eosio-enterprise/chappe/internal/encryption"
	"github.com/fatih/color"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/leekchan/accounting"
	"github.com/spf13/viper"
)

// LedgerTrx ...
type LedgerTrx struct {
	TrxDate  string `faker:"date"`
	Memo     string `faker:"sentence"`
	Account1 string `faker:"-"`
	Account2 string `faker:"-"`
	Amount   string `faker:"-"`
}

// Message ...
type Message struct {
	Payload map[string][]byte
}

// Bytes ...
func (m Message) Bytes() []byte {
	jsonMessage, err := json.Marshal(m)
	if err != nil {
		log.Println("Cannot convert to message to bytes", err)
	}

	return []byte(jsonMessage)
}

// NewMessage - creates a new Message object, provisions the Payload map
func NewMessage() Message {
	m := Message{}
	m.Payload = make(map[string][]byte)
	return m
}

// Load ...
func Load(channelName, ledgerFile, ipfsHash string) (Message, error) {
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

	var receivedLedgerTrx LedgerTrx
	err = json.Unmarshal(plaintext, &receivedLedgerTrx)

	textTrx := receivedLedgerTrx.TrxDate + "\t" + receivedLedgerTrx.Memo + "\n"
	textTrx = textTrx + "\t\t" + receivedLedgerTrx.Account1 + "\t\t" + receivedLedgerTrx.Amount + "\n"
	textTrx = textTrx + "\t\t" + receivedLedgerTrx.Account2 + "\n\n"

	log.Println("Adding new transaction to ledger: ")
	color.Cyan(textTrx)

	f, err := os.OpenFile(ledgerFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(textTrx)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

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

// Components ...
type Components struct {
	AccountIndex1 int `faker:"boundary_start=0, boundary_end=29"`
	AccountIndex2 int `faker:"boundary_start=0, boundary_end=29"`
	Dollars       int `faker:"boundary_start=10000, boundary_end=10000000"`
	Cents         int `faker:"boundary_start=0, boundary_end=100"`
}

// GetLedgerPayload ...
func GetLedgerPayload() LedgerTrx {
	var ledgerPayload LedgerTrx
	err := faker.FakeData(&ledgerPayload)
	if err != nil {
		log.Println("Cannot generate fake data: ", err)
	}

	var components Components
	err = faker.FakeData(&components)
	if err != nil {
		log.Println("Cannot generate fake data: ", err)
	}

	var accounts [30]string
	accounts[0] = "Assets:Checking Account"
	accounts[1] = "Assets:Reserve Account"
	accounts[2] = "Assets:Land"
	accounts[3] = "Equity:Builder"
	accounts[4] = "Equity:Fund Limited Partners"
	accounts[5] = "Equity:Private Equity Partners"
	accounts[6] = "Liabilities:Bank One Loan"
	accounts[7] = "Liabilities:PNC Line of Credit"
	accounts[8] = "Expenses:Financing:Appraisal Fee"
	accounts[9] = "Expenses:Financing:Mortgage Reg Tax"
	accounts[10] = "Expenses:Construction Contract"
	accounts[11] = "Expenses:Financing:Financing Fee"
	accounts[12] = "Expenses:Development:Development Fee"
	accounts[13] = "Expenses:Diligence:Title Insurance"
	accounts[14] = "Expenses:Builders Insurance"
	accounts[15] = "Expenses:Diligence:Traffic Study"
	accounts[16] = "Expenses:Diligence:Survey"
	accounts[17] = "Expenses:Diligence:Environmental"
	accounts[18] = "Expenses:Diligence:GeoTech"
	accounts[19] = "Expenses:Design Costs:Architectural & Engineering"
	accounts[20] = "Expenses:Finishing:Furniture"
	accounts[21] = "Expenses:Finishing:Flooring"
	accounts[22] = "Expenses:Finishing:Walls"
	accounts[23] = "Expenses:Office:Document Processing"
	accounts[24] = "Expenses:Office:Legal Fees"
	accounts[25] = "Expenses:Bank One Loan:Interest"
	accounts[26] = "Expenses:PNC Line of Credit:Interest"
	accounts[27] = "Income:Condo Sales"
	accounts[28] = "Income:Single Family Sales"
	accounts[29] = "Income:Rent"

	ledgerPayload.Account1 = accounts[components.AccountIndex1]
	ledgerPayload.Account2 = accounts[components.AccountIndex2]

	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	ledgerPayload.Amount = ac.FormatMoney(float64(components.Dollars) + (float64(components.Cents) / float64(100)))

	return ledgerPayload
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

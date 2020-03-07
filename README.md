
# Chappe Private Messaging
This project implements many-to-many encryption using ```eosio``` and IPFS. The use case is described in the [background document](BACKGROUND.md) (it's a bit outdated in terms of design though)

## Quick Start

Clone repo & build
``` bash
git clone https://github.com/eosio-enterprise/chappe
cd chappe
make  (or "go build")
```

Create config.yaml
``` bash
cat > config.yaml <<EOF
IPFS:
    Endpoint: localhost:5001
Eosio:
    Endpoint: https://kylin.eosusa.news
    PublishAccount: messengerbus
    PublishPrivateKey: 5KAP1zytghuvowgprSPLNasajibZcxf4KMgdgNbrNj98xhcGAUa
Dfuse:
    WSEndpoint: wss://kylin.eos.dfuse.io/v1/stream
    Origin: github.com/eosio-enterprise/chappe
    ApiKey: web_***  # Replace this, get one at dfuse.io
KeyDirectory: channels/
EOF
```

Create Channel
```
./chappe create key --channel-name chan4242
```

Run Subscriber
```
./chappe subscribe --channel-name chan4242
```

Open New Shell, and run Publisher
```
./chappe publish --channel-name chan4242 --readable-memo "This is human-readable, unencrypted memo"
```

## Menu
Run chappe
``` bash
➜ ./chappe
Welcome to Chappe Private Messaging for EOSIO

Usage:
  chappe [flags]
  chappe [command]

Available Commands:
  create      Create chappe artifacts
  get         Get chappe artifacts
  help        Help about any command
  publish     Publish a private message to a channel
  server      Run a server
  subscribe   Subscribe to a channel
  update      Print update instructions
  version     Print the version

Flags:
  -h, --help   help for chappe

Use "chappe [command] --help" for more information about a command.
```

**PRs Welcome**

## Usage
### Dependencies
#### IPFS
I run with go-ipfs:latest running in Docker. It does not work with Infura (header errors?).
``` bash
export ipfs_staging=</absolute/path/to/somewhere/>
export ipfs_data=</absolute/path/to/somewhere_else/>

docker run -d --name ipfs_host -v $ipfs_staging:/export -v $ipfs_data:/data/ipfs -p 4001:4001 -p 127.0.0.1:8080:8080 -p 127.0.0.1:5001:5001 ipfs/go-ipfs:latest
```

### Create a Key (Channel)
``` bash
./chappe create key --channel-name chan4242
```
The channel key is an assymetric RSA key. If you create a channel and you want another node to receive your messages, you would share the "chan4242.pem" file.

### Subscribe to the Channel
This runs a server, so fork your terminal shell to hold the ENV VARS intact. 
``` bash
./chappe subscribe --channel-name chan4242
```

(After I publish, I will describe what happens with the subscribe process)

### Publish to the Channel
On a separate tab, publish a message:
``` bash
./chappe publish --channel-name chan4242 --readable-memo "This is human-readable, unencrypted memo"
```

Currently, the publish command generates fake private data to be shared on the channel.

For example: 
``` json
{
  "RecordID": "94ffdee8-7c40-4133-a210-e57740cb7a99",
  "FirstName": "Cordie",
  "LastName": "Schinner",
  "DOB": "2013-05-11",
  "CreditCardNumber": "6011076284887079",
  "CreditCardType": "Discover",
  "Email": "hVNKeys@mCOjD.com",
  "TimeZone": "Australia/Perth",
  "AmountDue": "LYD 0.210000",
  "PhoneNumber": "110-537-9426",
  "SafeWord": "adipisci",
  "LastScan": "1973-09-25 17:56:48"
}
```

#### Hybrid Encryption
In order to support large messages (files), we use hybrid encryption as described in this paper (https://pdfs.semanticscholar.org/87ff/ea85fbf52e22e4808e1fcc9e40ead4ff7738.pdf). 

We generate a random symmetric key to encrypt the message, then use the channel's (recipient's) assymetric public key to encrypt the symmetric key. 

```chappe``` handles all of this but the purpose of this is for documentation of how it works
##### Step 1: Generate a Random AES Key (Symmetric)

Generate a one-time use key to encrypt the body of the message.
``` go
key := [32]byte{}
_, err := io.ReadFull(rand.Reader, key[:])
if err != nil {
    panic(err)
}
```
##### Step 2: Encrypt the Message Data with the AES Key  

``` go
block, err := aes.NewCipher(key[:])
if err != nil {
    return nil, err
}

gcm, err := cipher.NewGCM(block)
if err != nil {
    return nil, err
}

nonce := make([]byte, gcm.NonceSize())
_, err = io.ReadFull(rand.Reader, nonce)
if err != nil {
    return nil, err
}

aesEncryptedData := gcm.Seal(nonce, nonce, plaintext, nil), nil
```

##### Step 3: Encrypt the AES Key with the Channel's Private Key 
``` go
encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, key, label)
if err != nil {
    fmt.Fprintf(os.Stderr, "Error from RSA encryption: %s\n", err)
    return nil, err
}
```

##### Step 4: Publish Object to IPFS
The Encrypted Payload plus the Encrypted AES Key can be combined together, along with human-readable text to form the object. It is published to IPFS and a hash is returned.

``` go
type PersistedObject struct {
    EncryptedPayload   []byte
    EncryptedAESKey    []byte
    UnencryptedPayload string
}

jsonPayloadNode, err := json.Marshal(payload)
if err != nil {
    fmt.Fprintf(os.Stderr, "Could not marshal:  %s", err)
}

hash, err := sh.Add(strings.NewReader(string(jsonPayloadNode)))
if err != nil {
    fmt.Fprintf(os.Stderr, "Could not add data to IPFS: %s", err)
}
fmt.Println("IPFS Hash: ", hash)
```

##### Step 5: Publish IPFS Hash to EOSIO Blockchain
The user/node then publishes the IPFS hash to the appropriate blockchain, which records the event's existence, although elements of this metadata can be masked.
``` go
txOpts := &eos.TxOptions{}
if err := txOpts.FillFromChain(api); err != nil {
    panic(fmt.Errorf("filling tx opts: %s", err))
}

tx := eos.NewTransaction([]*eos.Action{message.NewPub(hash, readableMemo)}, txOpts)
_, packedTx, err := api.SignTransaction(tx, txOpts.ChainID, eos.CompressionNone)
if err != nil {
    panic(fmt.Errorf("sign transaction: %s", err))
}

response, err := api.PushTransaction(packedTx)
if err != nil {
    panic(fmt.Errorf("push transaction: %s", err))
}
```

### Back to Subscription

The inverse happens on the subscription side: 
- dfuse fires a websocket ( TODO: [ ] need to migrate to GraphQL)
- IPFS document is retrieved
- AES key is decrypted
- Message is decrypted


(TODO: need to add an signed acknowledgement back to the sender)
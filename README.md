
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
    Endpoint: ipfs.digscar.com:5001
Eosio:
    Endpoint: https://kylin.eosusa.news
    PublishAccount: messengerbus
    PublishPrivateKey: 5KAP1zytghuvowgprSPLNasajibZcxf4KMgdgNbrNj98xhcGAUa
Dfuse:
    Protocol: GraphQL
    GraphQLEndpoint: kylin.eos.dfuse.io:443
    Origin: github.com/eosio-enterprise/chappe
    ApiKey: web_***  # Replace this, get one at dfuse.io
KeyDirectory: channels/
PublishInterval: 10s  # Go Duration object
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
## Features
- Send/receive encrypted (or unencrypted) messages/documents using public or private EOSIO chains
- Messages are sent on channels, and all nodes with the channel key can read messages
- Optionally publish receipts (acknowledgements) signed with a node's key
- Support for large messages and files
- Optionally mask all metadata (publisher, type of message)
- (Coming Soon) Ability to set reveal parameters that automatically publish decrypted version after time elapses
- (Coming Soon) Hierarchies for data visibility

## Menu
Run chappe
``` bash
➜ ./chappe
Welcome to Chappe Private Messaging for EOSIO

Usage:
  chappe [flags]
  chappe [command]

Available Commands:
  create      Create chappe channel
  get         Get chappe message (via IPFS cid)
  help        Help about any command
  publish     Publish a private message to a channel
  subscribe   Subscribe to a channel
  version     Print the version

Flags:
  -h, --help   help for chappe

Use "chappe [command] --help" for more information about a command.
```

## Usage
### Configuration File
Chappe will locate a file named ```config.yaml``` by looking in the following folders: ```.```, ```configs```, ```/etc/chappe```, and ```$HOME/.chappe```. 

You can override any variable in the configuration file by setting an environment variable with a prefix of ```CHAPPE_```, followed by an all capital letter version of the variable that you want to override. 

### Create a Key (Channel)
``` bash
./chappe create key --channel-name chan4242
```
The channel key is an assymetric RSA key. If you create a channel and you want another node to receive your messages, you would share the "chan4242.pem" file.

### Subscribe to the Channel
This runs a server, so run it in a separate terminal.
``` bash
./chappe subscribe --channel-name chan4242
```

You can optionally request that the subscriber submit receipts/acknowledgements for each message. To prove that the receipient received and decrypted the message, the recipient's device key (unique to only that node) signs the decrypted message. This signature is posted to the blockchain, and the original sender may verify that the intended recipient(s) successfully received the message.

To send a receipt, pass the ```send-receipts``` or ```-r``` flag.
``` bash
./chappe subscribe --channel-name chan424 -r
```

### Publish to the Channel
On a separate tab, publish a message:
``` bash
./chappe publish --channel-name chan4242 --readable-memo "This is human-readable, unencrypted memo"
```

Currently, the publish command generates fake private data to be shared on the channel. More options will be added soon.
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

The payload data is encrypted (see below) and constructed into an object and saved to IPFS. Here's an example of what one of the objects looks like: http://ipfs.digscar.com:8080/ipfs/QmNLuCqYR23RLzkE8fZvrnhsfaYJiawWAXcs2miLdeckND

The blockchain transaction payload appears like this: 
``` json 
{
   "payload": [
      {
         "key": "cid",
         "value": "QmfBmT8CDSaRYQ6b7z1URNXZih7jWvgRdHw1oH5rQFAqoy"
      },
      {
         "key": "memo",
         "value": "foobars memo"
      }
   ]
}
```

#### Hybrid Encryption
In order to support large messages (files), we use hybrid encryption as described in this paper (https://pdfs.semanticscholar.org/87ff/ea85fbf52e22e4808e1fcc9e40ead4ff7738.pdf). 

We generate a random symmetric key for each message, then use the channel's (recipient's) assymetric public key to encrypt the symmetric key. 

```chappe``` handles all of this but the purpose of this is for documentation of how it works
1. Generate a Random AES Key (Symmetric)
2. Encrypt the Message Data with the AES Key
3. Encrypt the AES Key with the Channel's Private Key
4. Publish Message to IPFS
5. Publish IPFS CID (hash) to EOSIO Blockchain

### Dependencies
#### Dfuse
You'll need a dfuse API key. You can register for a free one at dfuse.io

#### IPFS
It's simple to run your own IPFS node using Docker with only these 3 commands, or you may use ```ipfs.digscar.com``` for light testing.
I run go-ipfs:latest running in Docker. 
``` bash
export ipfs_staging=</absolute/path/to/somewhere/>
export ipfs_data=</absolute/path/to/somewhere_else/>

docker run -d --name ipfs_host -v $ipfs_staging:/export -v $ipfs_data:/data/ipfs -p 4001:4001 -p 127.0.0.1:8080:8080 -p 127.0.0.1:5001:5001 ipfs/go-ipfs:latest
```
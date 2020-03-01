package message

import eos "github.com/eoscanada/eos-go"

func NewPub(ipfsHash, memo string) *eos.Action {
	return &eos.Action{
		Account: eos.AN("messengerbus"),
		Name:    eos.ActN("pub"),
		Authorization: []eos.PermissionLevel{
			{Actor: "messengerbus", Permission: eos.PN("active")},
		},
		ActionData: eos.NewActionData(Message{
			IpfsHash: ipfsHash,
			Memo:     memo,
		}),
	}
}

type Message struct {
	IpfsHash string `json:"ipfs_hash"`
	Memo     string `json:"memo"`
}

package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func generateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

func exportPrivateKey(privkey *rsa.PrivateKey) []byte {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return privkeyPem
}

func exportPublicKey(pubkey *rsa.PublicKey) ([]byte, error) {
	pubkeyBytes, err := x509.MarshalPKIXPublicKey(pubkey)
	if err != nil {
		return nil, err
	}
	pubkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubkeyBytes,
		},
	)
	return pubkeyPem, nil
}

func create(keyname string) (*rsa.PrivateKey, *rsa.PublicKey) {

	priv, pub := generateRsaKeyPair()

	// Export the keys to pem string
	privPem := exportPrivateKey(priv)
	pubPem, _ := exportPublicKey(pub)

	ioutil.WriteFile(viper.GetString("KeyDirectory")+keyname+".pub", pubPem, 0644)
	ioutil.WriteFile(viper.GetString("KeyDirectory")+keyname+".pem", privPem, 0644)

	return priv, pub
}

// MakeCreate ...
func MakeCreate() *cobra.Command {
	var command = &cobra.Command{
		Use:          "create",
		Short:        "Create a new private chappe channel",
		Example:      `  chappe create --channel-name chan649`,
		SilenceUsage: false,
	}

	var channelName string
	command.Flags().StringVarP(&channelName, "channel-name", "n", "", "Name of the private channel to create")

	command.RunE = func(command *cobra.Command, args []string) error {

		if len(channelName) == 0 {
			return fmt.Errorf("--channel-name required")
		}

		create(channelName)

		fmt.Println(
			`=======================================================================
key ` + channelName + ` created in files ` + channelName + `.pem (private) and ` + channelName + `.pub (public)
=======================================================================`)

		return nil
	}
	return command
}

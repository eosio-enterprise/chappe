package apps

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"

	"github.com/eosio-enterprise/chappe/pkg"
	"github.com/spf13/cobra"
)

const defaultVersion = "" // If we don't set version then we get latest

// GenerateRsaKeyPair ...
func GenerateRsaKeyPair() (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, _ := rsa.GenerateKey(rand.Reader, 4096)
	return privkey, &privkey.PublicKey
}

// ExportRsaPrivateKeyAsPemStr ...
func ExportRsaPrivateKeyAsPemStr(privkey *rsa.PrivateKey) []byte {
	privkeyBytes := x509.MarshalPKCS1PrivateKey(privkey)
	privkeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privkeyBytes,
		},
	)
	return privkeyPem
}

// ExportRsaPublicKeyAsPemStr ...
func ExportRsaPublicKeyAsPemStr(pubkey *rsa.PublicKey) ([]byte, error) {
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

	priv, pub := GenerateRsaKeyPair()

	// Export the keys to pem string
	privPem := ExportRsaPrivateKeyAsPemStr(priv)
	pubPem, _ := ExportRsaPublicKeyAsPemStr(pub)

	ioutil.WriteFile(keyname+".pub", pubPem, 0644)
	ioutil.WriteFile(keyname+".pem", privPem, 0644)

	return priv, pub
}

// MakeCreateKey ...
func MakeCreateKey() *cobra.Command {
	var keyCmd = &cobra.Command{
		Use:          "key",
		Short:        "Create a new key pair",
		Long:         `Creates a new RSA key pair stored with the keyname.pem and keyname.pub files.`,
		Example:      `  chappe create key --key-name MyKey`,
		SilenceUsage: true,
	}

	keyCmd.Flags().StringP("key-name", "n", "default", "key name")

	keyCmd.RunE = func(command *cobra.Command, args []string) error {
		keyName, _ := command.Flags().GetString("key-name")

		if len(keyName) == 0 {
			return fmt.Errorf("--key-name required")
		}

		create(keyName)

		fmt.Println(
			`=======================================================================
key ` + keyName + ` created.
=======================================================================
		
` + pkg.ThanksForUsing)

		return nil
	}

	return keyCmd
}

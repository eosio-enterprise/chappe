package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

// ParseRsaPublicKeyFromPem ...
func ParseRsaPublicKeyFromPem(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}

// ParseRsaPrivateKeyFromPem ...
func ParseRsaPrivateKeyFromPem(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func loadPublicKey(keyname string) *rsa.PublicKey {

	publicKeyPemStr, err := ioutil.ReadFile(viper.GetString("KeyDirectory") + keyname + ".pub")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from reading key file: %s\n", err)
		return nil
	}

	publicKey, err := ParseRsaPublicKeyFromPem(publicKeyPemStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from reading key file: %s\n", err)
		return nil
	}

	return publicKey
}

// RsaEncrypt ...
func RsaEncrypt(channelName string, payload []byte) ([]byte, error) {
	publicKey := loadPublicKey(channelName)
	label := []byte("chappe") // TODO: migrate to something else?

	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, payload, label)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from RSA encryption: %s\n", err)
		return nil, err
	}
	return encryptedData, nil
}

// RsaDecrypt ...
func RsaDecrypt(channelName string, payload []byte) ([]byte, error) {
	privateKey := load(channelName)
	label := []byte("chappe") // TODO: migrate to something else?

	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, payload, label)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "Cannot decrypt the AES key: %s\n", err)
		return nil, err
	}

	return plaintext, nil
}

func parseRsaPrivateKeyFromPem(privPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priv, nil
}

func load(keyname string) *rsa.PrivateKey {

	privateKeyPemStr, err := ioutil.ReadFile(viper.GetString("KeyDirectory") + keyname + ".pem")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from reading key file: %s\n", err)
		return nil
	}

	priv, err := parseRsaPrivateKeyFromPem(privateKeyPemStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error from reading key file: %s\n", err)
		return nil
	}
	return priv
}

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

// CreateChannel ...
func CreateChannel(keyname string) (*rsa.PrivateKey, *rsa.PublicKey) {

	priv, pub := generateRsaKeyPair()

	// Export the keys to pem string
	privPem := exportPrivateKey(priv)
	pubPem, _ := exportPublicKey(pub)

	ioutil.WriteFile(viper.GetString("KeyDirectory")+keyname+".pub", pubPem, 0644)
	ioutil.WriteFile(viper.GetString("KeyDirectory")+keyname+".pem", privPem, 0644)

	return priv, pub
}

package ssh

import (
	"github.houston.softwaregrp.net/Hercules/gotools/utils"
	"github.houston.softwaregrp.net/Hercules/gotools/utils/logging"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	sshKeyName = "arcsight-key"
	sshBitSize = 4096
)

var (
	keyPath               = filepath.FromSlash(filepath.Clean(filepath.Join(utils.GetHomeDirectory(), ".ssh")))
	PrivateKey, PublicKey = resolveKeyPair()
)

func resolveKeyPair() (string, string) {
	privateKeyPath := filepath.Join(keyPath, sshKeyName)
	publicKeyPath := privateKeyPath + ".pub"
	return privateKeyPath, publicKeyPath
}

type KeyPair struct {
	privateKey []byte
	publicKey  []byte
}

func (keyPair KeyPair) GetPrivateKey() []byte {
	return keyPair.privateKey
}

func (keyPair KeyPair) GetPublicKey() []byte {
	return keyPair.publicKey
}

func NewKeyPair(privateKey []byte, publicKey []byte) *KeyPair {
	return &KeyPair{privateKey, publicKey}
}

// Generates key pair
func GenerateKeyPair() *KeyPair {
	privateKey, err := generatePrivateKey(sshBitSize)
	if err != nil {
		logging.LogError(err.Error())
		os.Exit(1)
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		logging.LogError(err.Error())
		os.Exit(1)
	}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	return NewKeyPair(privateKeyBytes, publicKeyBytes)
}

// Propagates key pair to the local host
func SaveKeyPair(keyPair *KeyPair) {
	err := os.MkdirAll(keyPath, 0700)
	if err != nil {
		logging.LogError(err.Error())
		os.Exit(1)
	}

	err = writeKeyToFile(keyPair.privateKey, PrivateKey)
	if err != nil {
		logging.LogError(err.Error())
		os.Exit(1)
	}

	err = writeKeyToFile(keyPair.publicKey, PublicKey)
	if err != nil {
		logging.LogError(err.Error())
		os.Exit(1)
	}
}

// Propagates key pair to the remote host
func PropagateKeyPair(client *goph.Client, keyPair *KeyPair, privateKeyPath string, publicKeyPath string) {
	err := UploadBinaryFileFromMemory(client, privateKeyPath, &BinaryFile{Data: keyPair.GetPrivateKey(), Mode: 0600})
	if err != nil {
		logging.LogError(err.Error())
		os.Exit(1)
	}

	err = UploadBinaryFileFromMemory(client, publicKeyPath, &BinaryFile{Data: keyPair.GetPublicKey(), Mode: 0644})
	if err != nil {
		logging.LogError(err.Error())
		os.Exit(1)
	}
}

// Creates RSA private key of a specified byte size
func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// Encodes private key from RSA to PEM format
func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	ASN1DERKey := x509.MarshalPKCS1PrivateKey(privateKey)

	privatePEMBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   ASN1DERKey,
	}

	privatePEMKey := pem.EncodeToMemory(&privatePEMBlock)

	return privatePEMKey
}

// Takes rsa.PublicKey and return bytes suitable for writing to .pub file
func generatePublicKey(privateKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privateKey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}

// Writes key to a file
func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}
	return nil
}

// Takes binary private key in form of byte array and validates it
func IsPrivateKeyValid(bytes []byte) bool {
	_, err := ssh.ParseRawPrivateKey(bytes)
	if err != nil {
		logging.LogError(err.Error())
	}
	return err == nil
}

// Reads private key from file and validates it
func VerifyPrivateKeyFromFile(filename string) bool {
	var bytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return false
	}
	return IsPrivateKeyValid(bytes)
}

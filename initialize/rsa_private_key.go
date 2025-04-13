package initialize

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

var RSAPrivateKey *rsa.PrivateKey

func RsaDecryptInit(filePath string) (err error) {
	// Load private key from the given file path
	key, err := os.ReadFile(filePath)
	if err != nil {
		return errors.New("Failed to load private key (Error 1): " + err.Error())
	}

	// Decode the PEM encoded key
	block, _ := pem.Decode(key)
	if block == nil {
		return errors.New("Failed to load private key (Error 2):")
	}

	// Parse the PKCS1 private key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return errors.New("Failed to load private key (Error 3): " + err.Error())
	}

	// Store the private key in a global variable
	RSAPrivateKey = privateKey
	return err
}

func DecryptPassword(encryptedPassword string) ([]byte, error) {
	// Decode the base64 encoded encrypted password
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return nil, fmt.Errorf("Failed to decode ciphertext: %v", err)
	}

	// Decrypt the password using the RSA private key
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, RSAPrivateKey, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("Decryption failed: %v", err)
	}

	// Return the decrypted password
	return decrypted, nil
}

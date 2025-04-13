package test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// var password = "eDdka1FFVlR3enZ4ZVZLdVljQWt0RVNNT25mUHdpTU5tbXk0dkYzSThTcTUzRWN6VWl3STdIUzRTZTM1MXFROTl5V2xkbmtwOTQwZDVpZVl6b2NwZVF0RXNSc21aSmZ3a3RES3BwbVpWRURLNGJzZHVjSFhXTzd2eDY3VmFsQThjbjEwSnp2d0xNKzZVeHpiK2VnTTJqRUd6aFhTMGZEQ0ZmcEJPSEdmb1FMV1l5eTN3RWtZc2lFUzlxWjZ4WTlZbEN4Y2dibk9jeURuVFV0N3RlalM0UFMzR3BpMnFEWHRLWlFPVkpndEJqaTNWb1F2dG5yS3VpcURpSFhyaTdXVTRSY3BDbGcrb1UvLzcyc0FyN0huRkp1TjdWZHozSitmVFBWdWdiL0k2enhPQjhVVldsOUhxcit3UVkrZy9QckZZSWJ3RHVFSlBpVkpwbW5LUWROOUVRPT0="

var RSAPrivateKey *rsa.PrivateKey
var RSAPublicKey *rsa.PublicKey

func RsaDecryptInit(filePath string) (err error) {
	key, err := os.ReadFile(filePath)
	if err != nil {
		return errors.New("failed to read private key file: %v" + err.Error())
	}
	block, _ := pem.Decode(key)
	if block == nil {
		return errors.New("failed to decode PEM block from private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return errors.New("failed to parse PKCS#1 private key: %v" + err.Error())
	}
	RSAPrivateKey = privateKey

	return err
}

func RsaDecryptPublicInit(filePath string) (err error) {
	key, err := os.ReadFile(filePath)
	if err != nil {
		return errors.New("failed to read public key file: %v" + err.Error())
	}
	block, _ := pem.Decode(key)
	if block == nil {
		return errors.New("failed to decode PEM block from public key")
	}

	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return errors.New("failed to parse PKCS#1 public key: %v" + err.Error())
	}
	RSAPublicKey = publicKey
	return err
}

func DecryptPassword(encryptedPassword string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to decode encrypted password: %v", err)
	}

	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, RSAPrivateKey, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("RSA decryption failed: %v", err)
	}

	return decrypted, nil
}

func HashPassword(decryptedPassword []byte, _ []byte) (password []byte, err error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(decryptedPassword, bcrypt.DefaultCost)
	if err != nil {
		return password, fmt.Errorf("password hashing failed: %v", err)
	}
	return hashedPassword, err
}

func Encrypt() string {
	message := []byte("123456salt")
	encryptedMessage, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, RSAPublicKey, message, nil)
	if err != nil {
		log.Printf("RSA encryption failed: %v", err)
	}
	encryptPassword := base64.StdEncoding.EncodeToString(encryptedMessage)
	return encryptPassword
}

func TestRSA(t *testing.T) {
	// Initialize RSA private and public keys
	RsaDecryptInit("../rsa_key/private_key.pem")
	RsaDecryptPublicInit("../rsa_key/public.pem")

	// Perform encryption
	encryptedPassword := Encrypt()
	fmt.Println("Encrypted password:", encryptedPassword)

	// Decrypt the encrypted password
	decryptedBytes, err := DecryptPassword(encryptedPassword)
	if err != nil {
		return
	}

	// Output decrypted result
	fmt.Println("Decrypted value:", string(decryptedBytes))

	// Remove the fixed salt from the decrypted password
	password := strings.TrimRight(string(decryptedBytes), "salt")
	fmt.Println("Plain password:", password)

	t.Logf("Decrypted password (no salt): %v", password)
}

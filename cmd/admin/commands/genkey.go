package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/google/uuid"
)

func GenerateKey(keysize int) error {
	fmt.Printf("generating key with size %d...\n", keysize)

	privateKey, err := rsa.GenerateKey(rand.Reader, keysize)
	if err != nil {
		return fmt.Errorf("generate rsa private key: %w", err)
	}

	privateDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateDER,
	}

	//create folder keys if not already
	if err := os.Mkdir("keys", 0755); err != nil {
		if !os.IsExist(err) {
			return fmt.Errorf("creating keys folder: %w", err)
		}
	}

	//create the key file
	keyID := uuid.NewString()
	filepath := "keys/" + keyID + "-private.pem"

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating private key file: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, &privateBlock); err != nil {
		return fmt.Errorf("encoding into pem: %w", err)
	}
	// ==========================================================================
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshalling public key into DER: %w", err)
	}

	filepath = "keys/" + keyID + "-public.pem"

	publicFile, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("creating public key file: %w", err)
	}
	defer publicFile.Close()

	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyDER,
	}

	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return fmt.Errorf("encoding public key into pem: %w", err)
	}

	//make this key as active key
	activeKeyFilePath := "keys/active.txt"
	if err := os.WriteFile(activeKeyFilePath, []byte(keyID), 0644); err != nil {
		return fmt.Errorf("write active key file: %w", err)
	}

	fmt.Println("private and public key files generated")
	return nil
}

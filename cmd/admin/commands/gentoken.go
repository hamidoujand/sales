package commands

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const adminRole = "admin"

func GenerateToken(keyPath string, userId string, kid string) error {
	//TODO: need to add database check for the userId to make sure is authorized.
	claims := struct {
		jwt.RegisteredClaims
		Roles []string `json:"roles"`
	}{
		Roles: []string{adminRole},
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "admin-cli",
			Subject:   userId,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid

	//load the private key
	path := keyPath + "/" + kid + "-private.pem"

	privateFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open private key file: %w", err)
	}
	defer privateFile.Close()

	bs, err := io.ReadAll(privateFile)
	if err != nil {
		return fmt.Errorf("readAll: %w", err)
	}

	block, _ := pem.Decode(bs)
	if block == nil {
		//invalid pem data
		return fmt.Errorf("no pem block found in the file")
	}

	//parse the private key
	var privateKey *rsa.PrivateKey

	switch block.Type {
	case "RSA PRIVATE KEY":
		//use PKCS1
		var err error
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("parsing private key: %w", err)
		}

	case "PRIVATE KEY":
		key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("parsing private key: %w", err)
		}

		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return fmt.Errorf("key is not a rsa private key: %T", key)
		}

	default:
		return fmt.Errorf("unsupported PEM type: %s", block.Type)
	}

	//sign token
	tkn, err := token.SignedString(privateKey)
	if err != nil {
		return fmt.Errorf("signing token: %w", err)
	}

	fmt.Println("==============================TOKEN================================")
	fmt.Println(tkn)
	fmt.Println("==============================End================================")
	return nil
}

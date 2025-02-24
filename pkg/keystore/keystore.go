package keystore

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
)

var (
	ErrNotFound = errors.New("private key not found")
)

type KeyStore struct {
	store map[string]*rsa.PrivateKey
}

func New() *KeyStore {
	return &KeyStore{
		store: make(map[string]*rsa.PrivateKey),
	}
}

func (ks *KeyStore) LoadKeys(fsys fs.FS) (string, error) {

	//Example: c3550713-13e7-4a53-977a-dd53cbcb7088-private.pem
	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("trying to open %s: %w", path, err)
		}

		//skip dirs
		if d.IsDir() {
			return nil
		}

		//skip non-pem files
		if filepath.Ext(path) != ".pem" {
			return nil
		}

		//skip the public ones
		if strings.HasSuffix(path, "-public.pem") {
			return nil
		}

		//get the kid
		kid := strings.TrimSuffix(d.Name(), "-private.pem")

		file, err := fsys.Open(path)
		if err != nil {
			return fmt.Errorf("opening file %s: %w", path, err)
		}
		defer file.Close()

		//limit the read till 1MB
		bs, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			// 	return fmt.Errorf("reading key file %s: %w", path, err)
		}
		block, _ := pem.Decode(bs)
		if block == nil {
			return fmt.Errorf("invalid pem data")
		}

		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("parsing private key: %w", err)
		}
		ks.store[kid] = privateKey
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("walkdir: %w", err)
	}

	//find the active key
	activeKID, err := fsys.Open("active.txt")
	if err != nil {
		return "", fmt.Errorf("opening active kid file: %w", err)
	}
	kid, err := io.ReadAll(activeKID)
	if err != nil {
		return "", fmt.Errorf("readAll: %w", err)
	}

	return string(kid), nil
}

func (ks *KeyStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	key, ok := ks.store[kid]
	if !ok {
		return nil, ErrNotFound
	}

	return key, nil
}

func (ks *KeyStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	key, ok := ks.store[kid]
	if !ok {
		return nil, ErrNotFound
	}

	return &key.PublicKey, nil
}

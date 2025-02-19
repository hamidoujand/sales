package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hamidoujand/sales/internal/auth"
)

const kid = "key-id"

func TestAuth(t *testing.T) {
	issuer := "auth-service"
	s := newMockStore(t)
	a := auth.New(s, jwt.SigningMethodRS256, issuer)

	c := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   "user_id",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 2)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Roles: []string{"admin"},
	}
	token, err := a.GenerateToken(kid, c)
	if err != nil {
		t.Fatalf("failed to generate token: %s", err)
	}

	t.Logf("Token: %s\n", token)

	bearer := "Bearer " + token

	parsedClaims, err := a.Authenticate(context.Background(), bearer)
	if err != nil {
		t.Fatalf("failed to authenticate token: %s", err)
	}

	t.Logf("%+v\n", parsedClaims)
}

type mockStore struct {
	store map[string]*rsa.PrivateKey
}

func newMockStore(t *testing.T) *mockStore {
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %s", err)
	}

	s := mockStore{
		map[string]*rsa.PrivateKey{
			kid: private,
		},
	}
	return &s
}

func (ms *mockStore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	k := ms.store[kid]
	return k, nil
}

func (ms *mockStore) PublicKey(kid string) (*rsa.PublicKey, error) {
	k := ms.store[kid].PublicKey
	return &k, nil
}

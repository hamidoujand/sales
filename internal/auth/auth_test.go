package auth_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

	bearer := "Bearer " + token

	parsedClaims, err := a.Authenticate(context.Background(), bearer)
	if err != nil {
		t.Fatalf("failed to authenticate token: %s", err)
	}

	if parsedClaims.Issuer != c.Issuer {
		t.Errorf("issuer= %s, got %s", c.Issuer, parsedClaims.Issuer)
	}
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

func TestAuthorization(t *testing.T) {
	issuer := "auth-service"
	s := newMockStore(t)
	a := auth.New(s, jwt.SigningMethodRS256, issuer)

	tests := map[string]struct {
		claims     auth.Claims
		rule       string
		userId     string
		shouldFail bool
	}{
		"admin claims": {
			claims: auth.Claims{
				Roles: []string{"ADMIN"},
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer: issuer,
				},
			},
			rule:       auth.RuleAdmin,
			shouldFail: false,
			userId:     uuid.NewString(),
		},

		"user claim": {
			claims: auth.Claims{
				Roles: []string{"USER"},
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer: issuer,
				},
			},
			rule:       auth.RuleUser,
			userId:     uuid.NewString(),
			shouldFail: false,
		},

		"user accessing admin rule": {
			claims: auth.Claims{
				Roles: []string{"USER"},
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:  issuer,
					Subject: uuid.NewString(),
				},
			},
			rule:       auth.RuleAdmin,
			userId:     uuid.NewString(),
			shouldFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := a.Authorize(context.Background(), test.claims, test.userId, test.rule)
			if !test.shouldFail {
				if err != nil {
					t.Fatalf("failed to authorized with valid claims: %s", err)
				}
			} else {
				//we expected to fail
				if err == nil {
					t.Fatal("expected test to fail")
				}
			}
		})
	}
}

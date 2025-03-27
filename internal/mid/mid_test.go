package mid_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hamidoujand/sales/internal/auth"
	"github.com/hamidoujand/sales/internal/errs"
	"github.com/hamidoujand/sales/internal/mid"
	"github.com/hamidoujand/sales/internal/web"
)

const (
	roleAdmin = "ADMIN"
	roleUser  = "USER"
)

func TestAuthenticate(t *testing.T) {
	ks := newKeystroe(t)
	authClient := auth.New(ks, jwt.SigningMethodRS256, "auth-service", ks.activeKid)

	tests := map[string]struct {
		userId           string
		role             string
		expectErr        bool
		errStatus        int
		wrongTokenFormat bool
	}{
		"valid_bearer_token": {
			userId:           uuid.NewString(),
			role:             roleUser,
			expectErr:        false,
			wrongTokenFormat: false,
		},

		"invalid_bearer_format": {
			userId:           uuid.NewString(),
			role:             roleUser,
			expectErr:        true,
			errStatus:        http.StatusUnauthorized,
			wrongTokenFormat: true,
		},

		"invalid_userID_format": {
			userId:    "some-random-id-10102928192",
			role:      roleUser,
			expectErr: true,
			errStatus: http.StatusUnauthorized,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/v1/auth", nil)
			w := httptest.NewRecorder()
			c := auth.Claims{
				Roles: []string{test.role},
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   test.userId,
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}

			token, err := authClient.GenerateToken(c)
			if err != nil {
				t.Fatalf("failed to generate token: %s", err)
			}
			if test.wrongTokenFormat {
				r.Header.Set("Authorization", token)
			} else {
				r.Header.Set("Authorization", "Bearer "+token)
			}

			h := web.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				//check for userID
				fetchedID, err := auth.GetUserID(ctx)
				if err != nil {
					t.Fatalf("failed to get userID: %s", err)
				}
				if test.userId != fetchedID.String() {
					t.Errorf("userID=%s, got %s", test.userId, fetchedID)
				}
				w.WriteHeader(http.StatusOK)
				return nil
			})

			mid := mid.Authenticate(authClient)
			withAuth := mid(h)

			//call it
			err = withAuth(r.Context(), w, r)

			if !test.expectErr {
				if err != nil {
					t.Fatalf("failed to authenticate userbus: %s", err)
				}
			} else {
				if err == nil {
					t.Fatal("expected to authenticate to fail, but passed")
				}

				var appErr *errs.Error
				if !errors.As(err, &appErr) {
					t.Fatalf("expected the returned error to be of type errs.Error, got %T", err)
				}

				if test.errStatus != appErr.Code {
					t.Errorf("status=%d, got %d", test.errStatus, appErr.Code)
				}
			}
		})
	}
}

func TestAuthorize(t *testing.T) {
	ks := newKeystroe(t)
	authClient := auth.New(ks, jwt.SigningMethodRS256, "auth-service", ks.activeKid)
	tests := map[string]struct {
		userId                 uuid.UUID
		roles                  []string
		rules                  string
		expectErr              bool
		isEmptyUserIdAndClaims bool
		errStatus              int
	}{
		"success_path": {
			userId:    uuid.New(),
			roles:     []string{roleUser},
			rules:     auth.RuleUser,
			expectErr: false,
		},
		"no_userId_and_claims_in_ctx": {
			expectErr:              true,
			errStatus:              http.StatusUnauthorized,
			isEmptyUserIdAndClaims: true,
		},
		"user_can't_accessing_admin_route": {
			userId:    uuid.New(),
			roles:     []string{roleUser},
			rules:     auth.RuleAdmin,
			expectErr: true,
			errStatus: http.StatusUnauthorized,
		},

		"admin_can't_accessing_user_only_route": {
			userId:    uuid.New(),
			roles:     []string{roleAdmin},
			rules:     auth.RuleUser,
			expectErr: true,
			errStatus: http.StatusUnauthorized,
		},
		"admin_and_the_owner_can_access": {
			userId:    uuid.New(),
			roles:     []string{roleAdmin},
			rules:     auth.RuleAdminOrOwner,
			expectErr: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			c := auth.Claims{
				Roles: test.roles,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "auth-service",
					Subject:   test.userId.String(),
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			}
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/v1/test", nil)

			var ctx context.Context
			if !test.isEmptyUserIdAndClaims {
				ctx = r.Context()
				ctx = auth.SetUserId(ctx, test.userId)
				ctx = auth.SetClaims(ctx, c)
			} else {
				ctx = r.Context()
			}

			h := web.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				//check for userID
				fetchedID, err := auth.GetUserID(ctx)
				if err != nil {
					t.Fatalf("failed to get userID: %s", err)
				}
				if fetchedID != test.userId {
					t.Fatalf("userId=%s, got %s", test.userId, fetchedID)
				}
				w.WriteHeader(http.StatusOK)
				return nil
			})

			mid := mid.Authorize(authClient, test.rules)
			withAuth := mid(h)

			err := withAuth(ctx, w, r)
			if !test.expectErr {
				if err != nil {
					t.Fatalf("failed to authorize the request: %s", err)
				}
			} else {
				if err == nil {
					t.Fatal("expected to fail, but did not")
				}

				var appErr *errs.Error
				if !errors.As(err, &appErr) {
					t.Fatalf("expected error type to be errs.Error, got %T", err)
				}

				if test.errStatus != appErr.Code {
					t.Errorf("status=%d, got %d", test.errStatus, appErr.Code)
				}
			}
		})
	}

}

//==============================================================================

type keystore struct {
	store     map[string]*rsa.PrivateKey
	activeKid string
}

func newKeystroe(t *testing.T) *keystore {
	kid := uuid.NewString()
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal("failed to generate rsa private key")
	}

	return &keystore{
		store: map[string]*rsa.PrivateKey{
			kid: private,
		},
		activeKid: kid,
	}
}

func (ks *keystore) PrivateKey(kid string) (*rsa.PrivateKey, error) {
	return ks.store[kid], nil
}

func (ks *keystore) PublicKey(kid string) (*rsa.PublicKey, error) {
	return &ks.store[kid].PublicKey, nil
}

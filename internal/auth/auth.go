package auth

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"strings"
	"time"

	_ "embed"

	"github.com/golang-jwt/jwt/v5"
	"github.com/open-policy-agent/opa/v1/rego"
)

const regoPackageName = "token_validation"

var (
	//go:embed rego/authentication.rego
	regoAuthentication string
)

// KeyLookup defines the required behavior in order to get private and public keys for JWT token operations.
type KeyLookup interface {
	PrivateKey(kid string) (*rsa.PrivateKey, error)
	PublicKey(kid string) (*rsa.PublicKey, error)
}

type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles"`
}

type Auth struct {
	store         KeyLookup
	signingMethod jwt.SigningMethod
	issuer        string
}

func New(keyLookup KeyLookup, signingMethod jwt.SigningMethod, issuer string) *Auth {
	a := Auth{
		store:         keyLookup,
		signingMethod: signingMethod,
		issuer:        issuer,
	}
	return &a
}

// GenerateToken generates a jwt token based on the given claims.
func (a *Auth) GenerateToken(kid string, claims Claims) (string, error) {
	token := jwt.NewWithClaims(a.signingMethod, claims)
	token.Header["kid"] = kid

	//load the key
	privateKey, err := a.store.PrivateKey(kid)
	if err != nil {
		return "", fmt.Errorf("looking up private key: %w", err)
	}

	tkn, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return tkn, nil
}

func (a *Auth) Authenticate(ctx context.Context, bearerToken string) (Claims, error) {
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		return Claims{}, errors.New("expected Authorization header to be in this format: Bearer <TOKEN>")
	}

	tokenStr := strings.Split(bearerToken, " ")[1]

	var claims Claims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		rawKid, exists := t.Header["kid"]
		if !exists {
			return nil, fmt.Errorf("key id not found in the token header")
		}

		kid, ok := rawKid.(string)
		if !ok {
			return nil, fmt.Errorf("invalid key id")
		}

		//load the public key
		public, err := a.store.PublicKey(kid)
		if err != nil {
			return nil, fmt.Errorf("fetching public key: %w", err)
		}

		return public, nil
	})

	if err != nil {
		return Claims{}, fmt.Errorf("parsing token failed: %w", err)
	}

	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}

	//let the OPA to validate the claims
	input := map[string]any{
		"token": map[string]any{
			"iss":   a.issuer,
			"exp":   claims.ExpiresAt.Unix(),
			"roles": claims.Roles,
		},
		"now": time.Now().Unix(),
	}

	const validateRule = "valid"
	q := fmt.Sprintf("x = data.%s.%s", regoPackageName, validateRule)
	query, err := rego.New(
		rego.Query(q),
		rego.Module("policy.rego", regoAuthentication), // in case of any error they will shown like they are from a file named "policy.rego"
	).PrepareForEval(ctx)

	if err != nil {
		return Claims{}, fmt.Errorf("rego prepareForEval: %w", err)
	}

	results, err := query.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		return Claims{}, fmt.Errorf("query eval: %w", err)
	}

	if len(results) == 0 || len(results[0].Bindings) == 0 {
		return Claims{}, errors.New("no policy decision")
	}

	result, ok := results[0].Bindings["x"].(bool)
	if !result || !ok {
		return Claims{}, errors.New("access denied by policy")
	}
	return claims, nil
}

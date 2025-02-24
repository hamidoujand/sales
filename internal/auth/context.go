package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type ctxKey int

const userIdKey ctxKey = 1
const claimsKey ctxKey = 2

func SetUserId(ctx context.Context, userId uuid.UUID) context.Context {
	return context.WithValue(ctx, userIdKey, userId)
}

func SetClaims(ctx context.Context, claims Claims) context.Context {
	return context.WithValue(ctx, claimsKey, claims)
}

func GetUserID(ctx context.Context) (uuid.UUID, error) {
	userId, ok := ctx.Value(userIdKey).(uuid.UUID)
	if !ok {
		return uuid.Nil, fmt.Errorf("user id not found in the context")
	}

	return userId, nil
}

func GetClaims(ctx context.Context) (Claims, error) {
	c, ok := ctx.Value(claimsKey).(Claims)
	if !ok {
		return Claims{}, errors.New("claims not found in the context")
	}

	return c, nil
}

package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/hamidoujand/sales/internal/auth"
	"github.com/hamidoujand/sales/internal/errs"
	"github.com/hamidoujand/sales/internal/web"
)

func Authenticate(a *auth.Auth) web.Middleware {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			token := r.Header.Get("Authorization")
			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()

			claims, err := a.Authenticate(ctx, token)
			if err != nil {
				return errs.New(http.StatusUnauthorized, auth.ErrUnauthenticated)
			}

			//set claims into ctx
			userId, err := uuid.Parse(claims.Subject)
			if err != nil {
				return errs.Newf(http.StatusUnauthorized, "invalid userID: %s", userId)
			}
			ctx = auth.SetUserId(ctx, userId)
			ctx = auth.SetClaims(ctx, claims)
			return next(ctx, w, r)
		}
	}
}

package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/hamidoujand/sales/internal/auth"
	"github.com/hamidoujand/sales/internal/errs"
	"github.com/hamidoujand/sales/internal/web"
)

func Authorize(a *auth.Auth, rule string) web.Middleware {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			userId, err := auth.GetUserID(ctx)
			if err != nil {
				return errs.New(http.StatusUnauthorized, auth.ErrUnauthenticated)
			}

			ctx, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()

			claims, err := auth.GetClaims(ctx)
			if err != nil {
				return errs.New(http.StatusUnauthorized, auth.ErrUnauthenticated)
			}
			if err := a.Authorize(ctx, claims, userId.String(), rule); err != nil {
				return errs.New(http.StatusUnauthorized, auth.ErrUnauthenticated)
			}

			return next(ctx, w, r)
		}
	}
}

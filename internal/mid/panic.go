package mid

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/hamidoujand/sales/internal/errs"
	"github.com/hamidoujand/sales/internal/web"
)

func Panic() web.Middleware {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rc := recover(); rc != nil {
					//create a new trusted error from it
					trace := string(debug.Stack())
					err = errs.Newf(http.StatusInternalServerError, "PANIC[%v] TRACE[%s]", rc, trace)
				}
			}()

			return next(ctx, w, r)
		}
	}
}

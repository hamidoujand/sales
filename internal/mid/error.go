package mid

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/hamidoujand/sales/internal/errs"
	"github.com/hamidoujand/sales/internal/web"
)

func Error(log *slog.Logger) web.Middleware {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			//call the next handler
			err := next(ctx, w, r)
			if err == nil {
				return nil
			}

			//we have an err then

			var trustedErr *errs.Error
			if !errors.As(err, &trustedErr) {
				//untrusted error
				trustedErr = errs.Newf(http.StatusInternalServerError, "internal server error")
			}

			//trusted error
			log.Error("request failed",
				"code", trustedErr.Code,
				"funcName", filepath.Base(trustedErr.FuncName),
				"filename", filepath.Base(trustedErr.Filename),
				"msg", trustedErr.Message,
			)

			//internal server errors, not data should leak
			if trustedErr.Code == http.StatusInternalServerError {
				trustedErr.Message = http.StatusText(http.StatusInternalServerError)
			}

			if err := web.Respond(ctx, w, trustedErr.Code, trustedErr); err != nil {
				//log the error
				log.Error("responding error to client", "msg", err)
			}
			return nil
		}
	}
}

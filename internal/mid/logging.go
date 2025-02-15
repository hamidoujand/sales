package mid

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/hamidoujand/sales/internal/web"
)

func Logger(log *slog.Logger) web.Middleware {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			//before handler
			now := web.GetStartedAt(ctx)
			traceId := web.GetTraceID(ctx)

			path := r.URL.Path
			if r.URL.RawQuery != "" {
				path = fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
			}

			log.Info("request started", "traceID", traceId, "method", r.Method, "path", path, "remoteAddr", r.RemoteAddr)
			err := next(ctx, w, r)

			//after handler
			log.Info("request completed", "traceID", traceId, "method", r.Method, "statusCode", web.GetStatusCode(ctx), "path", path, "remoteAddr", r.RemoteAddr, "took", time.Since(now).String())
			return err
		}
	}
}

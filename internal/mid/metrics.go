package mid

import (
	"context"
	"net/http"

	"github.com/hamidoujand/sales/internal/metrics"
	"github.com/hamidoujand/sales/internal/web"
)

func Metrics() web.Middleware {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			//add metrics to the ctx
			ctx = metrics.Set(ctx)
			count := metrics.AddRequest(ctx)
			if count%1000 == 0 {
				metrics.AddGoroutines(ctx)
			}
			//call the handler
			err := next(ctx, w, r)
			if err != nil {
				metrics.AddError(ctx)
			}

			return err
		}
	}
}

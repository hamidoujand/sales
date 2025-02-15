package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/hamidoujand/sales/internal/mid"
	"github.com/hamidoujand/sales/internal/web"
)

func APIMux(logger *slog.Logger) *web.Router {
	const version = "v1"
	mux := web.NewRouter(logger, mid.Logger(logger))

	mux.HandleFunc(http.MethodGet, version, "/test/", testHandler)
	return mux
}

func testHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	msg := map[string]string{
		"msg": "Hello World!",
	}

	return web.Respond(ctx, w, http.StatusOK, msg)
}

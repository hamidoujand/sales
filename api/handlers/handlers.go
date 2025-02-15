package handlers

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"net/http"

	"github.com/hamidoujand/sales/internal/mid"
	"github.com/hamidoujand/sales/internal/web"
)

func APIMux(logger *slog.Logger) *web.Router {
	const version = "v1"
	mux := web.NewRouter(logger,
		mid.Logger(logger),
		mid.Error(logger),
		mid.Metrics(),
		mid.Panic(),
	)

	mux.HandleFunc(http.MethodGet, version, "/test/", testHandler)
	return mux
}

func testHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if rand.Int()%2 == 0 {
		//produce a dum error
		panic("something bad happened")
	}
	msg := map[string]string{
		"msg": "Hello World!",
	}

	return web.Respond(ctx, w, http.StatusOK, msg)
}

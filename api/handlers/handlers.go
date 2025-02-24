package handlers

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"net/http"

	"github.com/hamidoujand/sales/internal/auth"
	"github.com/hamidoujand/sales/internal/mid"
	"github.com/hamidoujand/sales/internal/web"
)

func APIMux(logger *slog.Logger, authClient *auth.Auth) *web.Router {
	const version = "v1"
	mux := web.NewRouter(logger,
		mid.Logger(logger),
		mid.Error(logger),
		mid.Metrics(),
		mid.Panic(),
	)

	mux.HandleFunc(http.MethodGet, version, "/test/", testHandler)
	mux.HandleFunc(http.MethodGet, version, "/authtest/", testAuthHandler, mid.Authenticate(authClient), mid.Authorize(authClient, auth.RuleAdmin))
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

func testAuthHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	msg := map[string]string{
		"msg": "Auth successful",
	}

	return web.Respond(ctx, w, http.StatusOK, msg)
}

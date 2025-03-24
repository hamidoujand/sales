package handlers

import (
	"log/slog"
	"net/http"

	"github.com/hamidoujand/sales/api/handlers/health"
	"github.com/hamidoujand/sales/internal/auth"
	"github.com/hamidoujand/sales/internal/mid"
	"github.com/hamidoujand/sales/internal/web"
	"github.com/jmoiron/sqlx"
)

func APIMux(build string, logger *slog.Logger, db *sqlx.DB, authClient *auth.Auth) *web.Router {
	const version = "v1"
	mux := web.NewRouter(logger,
		mid.Logger(logger),
		mid.Error(logger),
		mid.Metrics(),
		mid.Panic(),
	)

	//health handlers
	hh := health.Handler{
		DB:    db,
		Build: build,
	}

	mux.HandleFuncNoMid(http.MethodGet, version, "/readiness", hh.Readiness)
	mux.HandleFuncNoMid(http.MethodGet, version, "/liveness", hh.Liveness)

	return mux
}

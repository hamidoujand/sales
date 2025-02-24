package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hamidoujand/sales/api/handlers"
	"github.com/hamidoujand/sales/internal/auth"
	"github.com/hamidoujand/sales/internal/debug"
	"github.com/hamidoujand/sales/pkg/keystore"
)

var build = "development"

func main() {
	logger := configureLogger(os.Stdout, build, "sales")

	if err := run(logger); err != nil {
		logger.Error("run", "err", err.Error)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {

	//==========================================================================
	// Config
	cfg := struct {
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`   //TODO: needs load testing for actual value.
			IdleTimeout     time.Duration `conf:"default:120s"` //TODO: needs load testing for actual value.
			ShutdownTimeout time.Duration `conf:"default:20s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			APIHost         string        `conf:"default:0.0.0.0:8000"`
			DebugHost       string        `conf:"default:0.0.0.0:3000"`
		}

		Auth struct {
			KeysDir       string `conf:"default:keys"`
			SigningMethod string `conf:"default:RS256"`
			Issuer        string `conf:"default:auth-service"`
		}
	}{}

	help, err := conf.Parse("SALES", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parse config: %w", err)
	}

	//==========================================================================
	// GOMAXPROCS
	logger.Info("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	expvar.NewString("build").Set(build)

	confString, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("starting: %w", err)
	}

	logger.Info("startup", "configuration", confString)

	//==========================================================================
	// Debug Server
	go func() {
		logger.Info("debug server", "status", "running", "host", cfg.Web.DebugHost)
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux()); err != nil {
			logger.Error("debug server", "status", "failed", "err", err)
			return
		}
	}()

	//==========================================================================
	// Auth init
	ks := keystore.New()
	activeKid, err := ks.LoadKeys(os.DirFS(cfg.Auth.KeysDir))
	if err != nil {
		return fmt.Errorf("loading keys into key store: %w", err)
	}
	authClient := auth.New(ks, jwt.GetSigningMethod(cfg.Auth.SigningMethod), cfg.Auth.Issuer, activeKid)
	logger.Info("auth", "activeKID", activeKid)
	//==========================================================================
	// API server
	shutdown := make(chan os.Signal, 1)
	errCh := make(chan error, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	mux := handlers.APIMux(logger, authClient)

	server := &http.Server{
		Addr:        cfg.Web.APIHost,
		Handler:     http.TimeoutHandler(mux, cfg.Web.WriteTimeout, "time out"),
		ReadTimeout: cfg.Web.ReadTimeout,
		IdleTimeout: cfg.Web.IdleTimeout,
		ErrorLog:    slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	go func() {
		logger.Info("server started", "host", cfg.Web.APIHost)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("listen and serve: %w", err)
	case sig := <-shutdown:
		logger.Info("shutdown", "status", "shutting down", "signal", sig.String())
		defer logger.Info("shutdown", "status", "shutdown complete")

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("shutdown gracefully failed: %w", err)
		}
	}
	return nil
}

func configureLogger(w io.Writer, build string, service string) *slog.Logger {
	fn := func(groups []string, attr slog.Attr) slog.Attr {
		//customize the source attr
		if attr.Key == slog.SourceKey {
			if source, ok := attr.Value.Any().(*slog.Source); ok {
				filename := filepath.Base(source.File)
				filename = fmt.Sprintf("%s:%d", filename, source.Line)
				attr.Value = slog.StringValue(filename)
			}
		}
		return attr
	}

	var handler slog.Handler
	slogOpts := slog.HandlerOptions{AddSource: true, ReplaceAttr: fn}
	if build == "development" {
		handler = slog.NewTextHandler(w, &slogOpts)
	} else {
		handler = slog.NewJSONHandler(w, &slogOpts)
	}

	handler = handler.WithAttrs([]slog.Attr{
		{Key: "service", Value: slog.StringValue(service)},
		{Key: "build", Value: slog.StringValue(build)},
	})

	logger := slog.New(handler)
	return logger
}

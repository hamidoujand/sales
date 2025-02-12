package main

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
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
	// GOMAXPROCS
	logger.Info("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	//==========================================================================
	// Config
	cfg := struct {
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"`   //TODO: needs load testing for actual value.
			IdleTimeout     time.Duration `conf:"default:120s"` //TODO: needs load testing for actual value.
			ShutdownTimeout time.Duration `conf:"default:20s"`
			WriteTimeout    time.Duration `conf:"default:10s"`
			APIHost         string        `conf:"default:0.0.0.0:8000"`
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

	confString, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("string: %w", err)
	}

	logger.Info("startup", "configuration", confString)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown
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

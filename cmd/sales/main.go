package main

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
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

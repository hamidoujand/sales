package sqldb

import (
	"context"
	"fmt"
	"net/url"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host         string
	User         string
	Password     string
	Name         string
	Schema       string
	MaxIdleConns int
	MaxOpenConns int
	DisableTLS   bool
}

func Open(cfg Config) (*sqlx.DB, error) {
	q := make(url.Values)
	q.Set("timezone", "utc")

	sslMode := "required"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	q.Set("sslmode", sslMode)

	if cfg.Schema != "" {
		q.Set("search_path", cfg.Schema)
	}

	uri := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("pgx", uri.String())
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}

func StatusCheck(ctx context.Context, db *sqlx.DB) error {

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*10)
		defer cancel()
	}

	for attempts := 1; ; attempts++ {
		pingErr := db.PingContext(ctx)
		if pingErr == nil {
			break
		}

		time.Sleep(time.Duration(attempts) * 100 * time.Millisecond)
		//check ctx for another retry
		if ctx.Err() != nil {
			return fmt.Errorf("ctx deadline: %s", pingErr)
		}
	}

	//after ping, hit the sql engine
	var result bool
	if err := db.QueryRowContext(ctx, "SELECT true;").Scan(&result); err != nil {
		return fmt.Errorf("queryRowContext: %w", err)
	}

	return nil
}

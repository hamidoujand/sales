package sqldb

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const uniqueViolationCode = "23505"

var (
	ErrDuplicatedEntry = errors.New("duplicated entry")
)

//go:embed sql/*.sql
var migrationFiles embed.FS

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

func Migrate(ctx context.Context, db *sqlx.DB, dbname string) error {
	dirver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("creating postgres driver: %w", err)
	}

	src, err := iofs.New(migrationFiles, "sql") //prefix of the path: ie: "sql/init.sql"
	if err != nil {
		return fmt.Errorf("creating an iofs source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, dbname, dirver)
	if err != nil {
		return fmt.Errorf("creating a migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration up: %w", err)
	}

	return nil
}

func NamedExecContext(ctx context.Context, db *sqlx.DB, query string, data any) error {
	if _, err := db.NamedExecContext(ctx, query, data); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == uniqueViolationCode {
				return ErrDuplicatedEntry
			}
		}
		return fmt.Errorf("namedExecContext: %w", err)
	}
	return nil
}

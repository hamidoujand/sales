package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/hamidoujand/sales/internal/sqldb"
)

func Migrate(cfg sqldb.Config) error {
	fmt.Printf("applying migrations against %s.\n", cfg.Host)

	db, err := sqldb.Open(cfg)
	if err != nil {
		return fmt.Errorf("open: %w", err)
	}

	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	if err := sqldb.StatusCheck(ctx, db); err != nil {
		return fmt.Errorf("statusCheck: %w", err)
	}

	//run migrations
	if err := sqldb.Migrate(ctx, db, cfg.Name); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	fmt.Println("migrations applied successfully.")
	return nil
}

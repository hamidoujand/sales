// Package dbtest provides setup functionality for tests involved with database.
package dbtest

import (
	"context"
	"fmt"
	"math/rand/v2"

	"testing"

	"github.com/hamidoujand/sales/internal/sqldb"
	"github.com/hamidoujand/sales/pkg/docker"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	DB *sqlx.DB
}

func NewDatabase(ctx context.Context, t *testing.T, containerName string) *Database {
	image := "postgres:17.2"
	port := "5432"
	dockerArgs := []string{"-e", "POSTGRES_PASSWORD=password"}
	appArgs := []string{"-c", "log_statement=all"}

	c, err := docker.StartContainer(image, containerName, port, dockerArgs, appArgs)
	if err != nil {
		t.Fatalf("startContainer: %s", err)
	}

	t.Logf("Name    : %s\n", c.Name)
	t.Logf("HostPort: %s\n", c.HostPort)

	dbMaster, err := sqldb.Open(sqldb.Config{
		Host:       c.HostPort,
		User:       "postgres",
		Password:   "password",
		Name:       "postgres",
		DisableTLS: true,
	})

	if err != nil {
		t.Fatalf("open db conn: %s", err)
	}

	if err := sqldb.StatusCheck(ctx, dbMaster); err != nil {
		t.Fatalf("statusCheck: %s", err)
	}

	//--------------------------------------------------------------------------

	//random db name, must start with letter otherwise postgres complains.
	letters := "abcdefghijklmnopqrstuvwxyz"
	bs := make([]byte, 4)

	for i := range bs {
		idx := rand.IntN(len(letters))
		bs[i] = letters[idx]
	}

	dbname := string(bs)

	t.Logf("creating database: %s\n", dbname)
	if _, err := dbMaster.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", dbname)); err != nil {
		t.Fatalf("creating database %s: %s", dbname, err)
	}

	//connect to it
	db, err := sqldb.Open(sqldb.Config{
		Host:       c.HostPort,
		User:       "postgres",
		Password:   "password",
		Name:       dbname,
		DisableTLS: true,
	})

	if err != nil {
		t.Fatalf("connecting to db %s: %s", dbname, err)
	}

	t.Logf("running migrations against db %s\n", dbname)
	if err := sqldb.Migrate(ctx, db, dbname); err != nil {
		t.Fatalf("migration failed againt db %s: %s", dbname, err)
	}

	//register clean up
	t.Cleanup(func() {
		t.Helper()

		//close the db
		if err := db.Close(); err != nil {
			t.Fatalf("closing database %s", dbname)
		}

		// Terminate all connections to the database
		terminateSQL := `
		  SELECT pg_terminate_backend(pg_stat_activity.pid)
		  FROM pg_stat_activity
		  WHERE pg_stat_activity.datname = $1
		  AND pid <> pg_backend_pid()`

		if _, err := dbMaster.Exec(terminateSQL, dbname); err != nil {
			t.Errorf("terminating connections to %s: %v", dbname, err)
			return
		}

		//drop the database
		t.Logf("drop database %s", dbname)

		cleanUpCtx := context.Background() //for clean up use a background ctx so cleaup finishes.
		if _, err := dbMaster.ExecContext(cleanUpCtx, fmt.Sprintf("DROP DATABASE %s", dbname)); err != nil {
			t.Fatalf("dropping database %s: %s", dbname, err)
		}

		if err := dbMaster.Close(); err != nil {
			t.Fatalf("closing master db conn: %s", err)
		}

		//stop the container
		if err := docker.StopContainer(c.Name); err != nil {
			t.Fatalf("stopping container: %s", err)
		}
	})

	return &Database{
		DB: db,
	}
}

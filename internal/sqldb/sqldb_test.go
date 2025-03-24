package sqldb

import (
	"context"
	"testing"
	"time"

	"github.com/hamidoujand/sales/pkg/docker"
)

func TestDatabaseConn(t *testing.T) {
	image := "postgres:17.2"
	containerName := "test_postgres_conn"
	port := "5432"
	dockerArgs := []string{"-e", "POSTGRES_PASSWORD=password", "-e", "POSTGRES_USER=postgres", "-e", "POSTGRES_DB=postgres"}
	containerArgs := []string{"-c", "log_statement=all"}

	c, err := docker.StartContainer(image, containerName, port, dockerArgs, containerArgs)
	if err != nil {
		t.Fatalf("failed to start postgres container: %s", err)
	}

	defer func() { _ = docker.StopContainer(c.Name) }()
	cfg := Config{
		Host:       c.HostPort,
		User:       "postgres",
		Password:   "password",
		Name:       "postgres",
		DisableTLS: true,
	}

	db, err := Open(cfg)
	if err != nil {
		t.Fatalf("failed to open a conn: %s", err)
	}
	defer func() { _ = db.Close() }()

	//status
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	if err := StatusCheck(ctx, db); err != nil {
		t.Fatalf("statusCheck failed: %s", err)
	}
}

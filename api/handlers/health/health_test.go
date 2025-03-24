package health_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hamidoujand/sales/api/handlers/health"
	"github.com/hamidoujand/sales/internal/sqldb"
	"github.com/hamidoujand/sales/pkg/docker"
)

func TestHandler_Liveness(t *testing.T) {
	build := "TEST"
	image := "postgres:17.2"
	args := []string{"-e", "POSTGRES_PASSWORD=password", "-e", "POSTGRES_DB=postgres", "-e", "POSTGRES_USER=postgres"}
	containerArgs := []string{"-c", "log_statement=all"}
	c, err := docker.StartContainer(image, "health_check_liveness", "5432", args, containerArgs)
	if err != nil {
		t.Fatalf("failed to start the container: %s", err)
	}

	defer func() {
		_ = docker.StopContainer(c.Name)
	}()

	db, err := sqldb.Open(sqldb.Config{
		Host:       c.HostPort,
		User:       "postgres",
		Password:   "password",
		Name:       "postgres",
		DisableTLS: true,
	})

	if err != nil {
		t.Fatalf("failed to open a db conn: %s", err)
	}

	defer func() {
		_ = db.Close()
	}()

	hh := health.Handler{
		DB:    db,
		Build: build,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/liveness", nil)
	w := httptest.NewRecorder()

	err = hh.Liveness(context.Background(), w, req)
	if err != nil {
		t.Fatalf("failed to call liveness: %s", err)
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("status=%d, got=%d", http.StatusOK, w.Result().StatusCode)
	}

	var resp health.Info
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal into health.Info: %s", err)
	}

	if resp.Build != build {
		t.Errorf("build=%s, got=%s", build, resp.Build)
	}

	if resp.Status != "up" {
		t.Errorf("status=%s, got=%s", "up", resp.Status)
	}
}

func TestHandler_Readiness(t *testing.T) {
	build := "TEST"
	image := "postgres:17.2"
	args := []string{"-e", "POSTGRES_PASSWORD=password", "-e", "POSTGRES_DB=postgres", "-e", "POSTGRES_USER=postgres"}
	containerArgs := []string{"-c", "log_statement=all"}
	c, err := docker.StartContainer(image, "health_check_liveness", "5432", args, containerArgs)
	if err != nil {
		t.Fatalf("failed to start the container: %s", err)
	}

	defer func() {
		_ = docker.StopContainer(c.Name)
	}()

	db, err := sqldb.Open(sqldb.Config{
		Host:       c.HostPort,
		User:       "postgres",
		Password:   "password",
		Name:       "postgres",
		DisableTLS: true,
	})

	if err != nil {
		t.Fatalf("failed to open a db conn: %s", err)
	}

	defer func() {
		_ = db.Close()
	}()

	hh := health.Handler{
		DB:    db,
		Build: build,
	}

	req := httptest.NewRequest(http.MethodGet, "/v1/readiness", nil)
	w := httptest.NewRecorder()

	err = hh.Readiness(context.Background(), w, req)
	if err != nil {
		t.Fatalf("failed to call readiness: %s", err)
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Fatalf("status=%d, got=%d", http.StatusOK, w.Result().StatusCode)
	}
}

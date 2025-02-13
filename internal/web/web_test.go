package web_test

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hamidoujand/sales/internal/web"
)

type data struct {
	counter int
}

type ctxKey int

const dataKey ctxKey = 1

func TestRouter(t *testing.T) {
	log := slog.New(slog.DiscardHandler)
	r := web.NewRouter(log, mid1)

	server := httptest.NewServer(r)
	defer server.Close()

	r.HandleFunc(http.MethodGet, "v1", "/test", handler(t), mid2, mid3)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/v1/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("failed to make the request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func mid1(next web.HandlerFunc) web.HandlerFunc {
	//inject data into ctx
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		d := data{counter: 0}
		ctx = context.WithValue(ctx, dataKey, &d)
		return next(ctx, w, r)
	}
}

func mid2(next web.HandlerFunc) web.HandlerFunc {
	//increment
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		d := ctx.Value(dataKey).(*data)
		d.counter++
		ctx = context.WithValue(ctx, dataKey, d)
		return next(ctx, w, r)
	}
}

func mid3(next web.HandlerFunc) web.HandlerFunc {
	//increment
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		d := ctx.Value(dataKey).(*data)
		d.counter++
		ctx = context.WithValue(r.Context(), dataKey, d)
		return next(ctx, w, r)
	}
}
func handler(t *testing.T) web.HandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		d := ctx.Value(dataKey).(*data)
		var hasFailed bool

		if d.counter != 2 {
			hasFailed = true
			t.Errorf("counter=%d, got %d", 2, d.counter)
		}

		if r.Method != http.MethodGet {
			hasFailed = true
			t.Errorf("method=%s, got %s", http.MethodGet, r.Method)
		}
		if r.URL.Path != "/v1/test" {
			hasFailed = true
			t.Errorf("path=%s, got %s", "/v1/test", r.URL.Path)
		}
		statusCode := http.StatusOK

		if hasFailed {
			statusCode = http.StatusBadRequest
		}
		w.WriteHeader(statusCode)
		return nil
	}
}

func TestRespond(t *testing.T) {
	w := httptest.NewRecorder()
	ctx := context.Background()
	msg := "hello world!"
	data := map[string]string{
		"msg": msg,
	}
	statusCode := http.StatusOK

	if err := web.Respond(ctx, w, statusCode, data); err != nil {
		t.Fatalf("failed to respond: %s", err)
	}

	if w.Result().StatusCode != statusCode {
		t.Errorf("status=%d, got %d", w.Result().StatusCode, statusCode)
	}

	var result map[string]string
	bs, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("failed to read all response body: %s", err)
	}

	if err := json.Unmarshal(bs, &result); err != nil {
		t.Fatalf("failed to unmarshal data: %s", err)
	}

	received := result["msg"]
	if msg != received {
		t.Fatalf("msg=%s, got %s", msg, received)
	}
}

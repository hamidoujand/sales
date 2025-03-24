package health

import (
	"context"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/hamidoujand/sales/internal/sqldb"
	"github.com/hamidoujand/sales/internal/web"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	DB    *sqlx.DB
	Build string
}

type Info struct {
	Status     string `json:"status,omitempty"`
	Build      string `json:"build,omitempty"`
	Host       string `json:"host,omitempty"`
	Name       string `json:"name,omitempty"`
	PodIP      string `json:"podIP,omitempty"`
	Node       string `json:"node,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	GOMAXPROCS int    `json:"GOMAXPROCS,omitempty"`
}

func (h *Handler) Liveness(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	info := Info{
		Status:     "up",
		Build:      h.Build,
		Host:       host,
		Name:       os.Getenv("KUBERNETES_NAME"),
		PodIP:      os.Getenv("KUBERNETES_POD_IP"),
		Node:       os.Getenv("KUBERNETES_NODE_NAME"),
		Namespace:  os.Getenv("KUBERNETES_NAMESPACE"),
		GOMAXPROCS: runtime.GOMAXPROCS(0),
	}

	return web.Respond(ctx, w, http.StatusOK, info)
}

func (h *Handler) Readiness(ctx context.Context, w http.ResponseWriter, _ *http.Request) error {
	ctx, cancel := context.WithTimeout(ctx, time.Minute) //slow machine
	defer cancel()

	if err := sqldb.StatusCheck(ctx, h.DB); err != nil {
		return web.Respond(ctx, w, http.StatusInternalServerError, nil)
	}
	return web.Respond(ctx, w, http.StatusOK, nil)
}

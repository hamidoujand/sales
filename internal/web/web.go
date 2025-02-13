package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

type HandlerFunc func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

type Router struct {
	*http.ServeMux
	log  *slog.Logger
	mids []Middleware //global middlewares
}

func NewRouter(logger *slog.Logger, mids ...Middleware) *Router {
	return &Router{
		mids:     mids,
		log:      logger,
		ServeMux: http.NewServeMux(),
	}
}

func (r *Router) HandleFunc(method string, version string, path string, handlerFunc HandlerFunc, mids ...Middleware) {
	handler := applyMiddleware(handlerFunc, mids...)
	handler = applyMiddleware(handler, r.mids...)

	h := func(w http.ResponseWriter, req *http.Request) {
		//this is the actual outer layer that will be called by serveMux , here is where
		//we call our own custom handler.

		ctx := req.Context()
		if err := handler(ctx, w, req); err != nil {
			//with proper error handler middleware, we should not get an error in here
			//if it did, we just log it.
			r.log.Error("router", "status", "handlerFunc", "err", err)
			return
		}
	}

	pattern := path
	if version != "" {
		pattern = "/" + version + path
	}

	pattern = fmt.Sprintf("%s %s", method, pattern)
	r.ServeMux.HandleFunc(pattern, h)
}

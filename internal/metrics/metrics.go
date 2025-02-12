package metrics

import (
	"context"
	"expvar"
	"runtime"
)

// metrics is going to be used as singleton, will be injected into request ctx
// and each middleware can inc those counters.
var m *metrics

type metrics struct {
	goroutines *expvar.Int
	requests   *expvar.Int
	panics     *expvar.Int
	errors     *expvar.Int
}

func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests:   expvar.NewInt("requests"),
		panics:     expvar.NewInt("panics"),
		errors:     expvar.NewInt("errors"),
	}
}

type ctxKey int

const key ctxKey = 1

func Set(ctx context.Context) context.Context {
	return context.WithValue(ctx, key, m)
}

func AddRequest(ctx context.Context) int64 {
	m, ok := ctx.Value(key).(*metrics)
	if ok {
		m.requests.Add(1)
		return m.requests.Value()
	}
	return 0
}

func AddPanic(ctx context.Context) int64 {
	m, ok := ctx.Value(key).(*metrics)
	if ok {
		m.panics.Add(1)
		return m.panics.Value()
	}
	return 0
}

func AddError(ctx context.Context) int64 {
	m, ok := ctx.Value(key).(*metrics)
	if ok {
		m.errors.Add(1)
		return m.errors.Value()
	}
	return 0
}

func AddGoroutines(ctx context.Context) int64 {
	m, ok := ctx.Value(key).(*metrics)
	if ok {
		g := runtime.NumGoroutine()
		m.goroutines.Set(int64(g))
		return m.goroutines.Value()
	}
	return 0
}

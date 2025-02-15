package web

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// RequestData represents additional data we want to collect arount each request.
type RequestData struct {
	//TODO: traceID should be able to pass it cross processes in microservices, so will be replaced with openTel traceId.
	TraceID    string
	StartedAt  time.Time
	StatusCode int
}

type ctxKey int

const reqKey ctxKey = 1

func setRequestData(ctx context.Context) context.Context {
	rd := RequestData{
		TraceID:   uuid.NewString(),
		StartedAt: time.Now(),
	}
	return context.WithValue(ctx, reqKey, &rd)
}

func GetTraceID(ctx context.Context) string {
	rd, ok := ctx.Value(reqKey).(*RequestData)
	if !ok {
		return uuid.Nil.String()
	}

	return rd.TraceID
}

func GetStartedAt(ctx context.Context) time.Time {
	rd, ok := ctx.Value(reqKey).(*RequestData)
	if !ok {
		return time.Now()
	}
	return rd.StartedAt
}

func setStatusCode(ctx context.Context, statusCode int) bool {
	rd, ok := ctx.Value(reqKey).(*RequestData)
	if !ok {
		return false
	}
	rd.StatusCode = statusCode
	return true
}

func GetStatusCode(ctx context.Context) int {
	rd, ok := ctx.Value(reqKey).(*RequestData)
	if !ok {
		return 0
	}

	return rd.StatusCode
}

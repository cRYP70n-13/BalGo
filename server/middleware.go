package server

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Middleware struct {
	mux http.Handler
}

type contextValueKey string

var (
	user             contextValueKey = "user"
	requestStartTime contextValueKey = "__request_start_time__"
	cid              contextValueKey = "cid"
)

// TODO: Improve this middleware as well and add another one like an interceptor
// to add metadata to the outgoing requests.
func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := context.WithValue(r.Context(), user, "otmane_kimdil")
	ctx = context.WithValue(ctx, requestStartTime, time.Now())
	ctx = context.WithValue(ctx, cid, uuid.New())
	req := r.WithContext(ctx)

	m.mux.ServeHTTP(w, req)

	start, ok := req.Context().Value(requestStartTime).(time.Time)
	if !ok {
		log.Fatalln("cannot get data from context conversion went wrong")
	}
	slog.Info("Request",
		"origin", r.RemoteAddr,
		"method", r.Method,
		"url", r.URL,
		"host", r.Host,
		"user-agent", r.UserAgent(),
        "request ID: ", req.Context().Value(cid),
        "request duration: ", time.Since(start),
	)
}

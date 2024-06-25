package server

import (
	"context"
	"log"
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
	log.Println("request duration: ", time.Since(start))
	log.Println("request ID: ", req.Context().Value(cid))
}

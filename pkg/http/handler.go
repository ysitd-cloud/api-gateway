package http

import (
	"net/http"
	"strings"

	"app.ysitd/gateway/pkg/http/grpc"
	"app.ysitd/gateway/pkg/http/rest"
	"context"
	"time"
)

const requestTimeout = 1 * time.Minute

type Handler struct {
	RestHandler *rest.Mux    `inject:""`
	GrpcHandler *grpc.Server `inject:""`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	r = r.WithContext(ctx)
	if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
		h.GrpcHandler.ServeHTTP(w, r)
	} else {
		h.RestHandler.ServeHTTP(w, r)
	}
}

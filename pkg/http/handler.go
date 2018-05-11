package http

import (
	"net/http"
	"strings"

	"app.ysitd/gateway/pkg/http/grpc"
	"app.ysitd/gateway/pkg/http/rest"
)

type Handler struct {
	RestHandler *rest.Mux    `inject:""`
	GrpcHandler *grpc.Server `inject:""`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
		h.RestHandler.ServeHTTP(w, r)
	} else {
		h.RestHandler.ServeHTTP(w, r)
	}
}

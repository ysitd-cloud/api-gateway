package grpc

import (
	"net/http"

	"google.golang.org/grpc"

	"golang.ysitd.cloud/log"

	"code.ysitd.cloud/api/account"
)

type Server struct {
	Server *grpc.Server `inject:""`
	Logger log.Logger   `inject:"grpc logger"`

	AccountBackend *account.Client `inject:""`
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Server.ServeHTTP(w, r)
}

package grpc

import (
	"net/http"
	"sync"

	"google.golang.org/grpc"

	"golang.ysitd.cloud/log"

	"code.ysitd.cloud/api/account"
	"code.ysitd.cloud/api/totp"
)

type Server struct {
	bootstrap sync.Once

	Server *grpc.Server `inject:""`
	Logger log.Logger   `inject:"grpc logger"`

	AccountBackend *account.Client `inject:""`
	TotpBackend    *totp.Client    `inject:""`
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.bootstrap.Do(func() {
		account.RegisterAccountServer(s.Server, s)
		totp.RegisterTotpServer(s.Server, s)
	})
	s.Server.ServeHTTP(w, r)
}

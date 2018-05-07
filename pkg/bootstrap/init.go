package bootstrap

import (
	"net/http"

	"code.ysitd.cloud/gateway/pkg/grpc"
	proxy "code.ysitd.cloud/gateway/pkg/http"
	"github.com/facebookgo/inject"
)

var grpcHandler *grpc.Server
var httpHandler *proxy.Mux

func init() {
	var graph inject.Graph
	graph.Logger = initLogger()

	graph.Provide(
		&inject.Object{Value: &grpcHandler},
		&inject.Object{Value: &httpHandler},
	)

	for _, fn := range []func(*inject.Graph){
		injectLogger,
		injectGrpc,
		injectBackend,
	} {
		fn(&graph)
	}

	if err := graph.Populate(); err != nil {
		panic(err)
	}
}

func GetGrpcHandler() http.Handler {
	return grpcHandler
}

func GetHttpHandler() http.Handler {
	return httpHandler
}

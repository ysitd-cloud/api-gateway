package bootstrap

import (
	"net/http"

	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"

	"code.ysitd.cloud/gateway/pkg/grpc"
	proxy "code.ysitd.cloud/gateway/pkg/http"
)

var grpcHandler grpc.Server
var httpHandler proxy.Mux

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
	return &grpcHandler
}

func GetHttpHandler() http.Handler {
	return &httpHandler
}

func GetMainLogger() logrus.FieldLogger {
	return logger.WithField("source", "main")
}

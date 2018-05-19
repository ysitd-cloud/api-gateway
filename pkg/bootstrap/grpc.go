package bootstrap

import (
	"github.com/facebookgo/inject"
	"google.golang.org/grpc"
)

func injectGrpc(graph *inject.Graph) error {
	return graph.Provide(
		&inject.Object{Value: grpc.NewServer()},
	)
}

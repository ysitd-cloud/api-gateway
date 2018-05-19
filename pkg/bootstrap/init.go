package bootstrap

import (
	"net/http"

	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"

	proxy "app.ysitd/gateway/pkg/http"
)

var handler proxy.Handler

type injectFunc func(*inject.Graph) error

func init() {
	var graph inject.Graph
	graph.Logger = initLogger()

	graph.Provide(
		&inject.Object{Value: &handler},
	)

	for _, fn := range []injectFunc{
		injectLogger,
		injectGrpc,
		injectBackend,
	} {
		if err := fn(&graph); err != nil {
			panic(err)
		}
	}

	if err := graph.Populate(); err != nil {
		panic(err)
	}
}

func GetHandler() http.Handler {
	return &handler
}

func GetMainLogger() logrus.FieldLogger {
	return logger.WithField("source", "main")
}

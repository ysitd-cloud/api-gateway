package bootstrap

import (
	"net/http"

	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"

	proxy "app.ysitd/gateway/pkg/http"
)

var handler proxy.Handler

func init() {
	var graph inject.Graph
	graph.Logger = initLogger()

	graph.Provide(
		&inject.Object{Value: &handler},
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

func GetHandler() http.Handler {
	return &handler
}

func GetMainLogger() logrus.FieldLogger {
	return logger.WithField("source", "main")
}

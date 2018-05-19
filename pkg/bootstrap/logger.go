package bootstrap

import (
	"os"

	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func initLogger() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()
		if os.Getenv("VERBOSE") != "" {
			logger.SetLevel(logrus.DebugLevel)
		}
	}

	return logger
}

func injectLogger(graph *inject.Graph) error {
	logger := initLogger()
	return graph.Provide(
		&inject.Object{Value: logger},
		&inject.Object{Value: logger.WithField("source", "http"), Name: "http logger"},
		&inject.Object{Value: logger.WithField("source", "grpc"), Name: "grpc logger"},
		&inject.Object{Value: logger.WithField("source", "account_proxy"), Name: "account proxy logger"},
	)
}

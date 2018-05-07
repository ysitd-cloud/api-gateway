package bootstrap

import (
	"os"

	"github.com/facebookgo/inject"

	accountApi "code.ysitd.cloud/api/account"
	accountClient "code.ysitd.cloud/client/account"
)

func injectBackend(graph *inject.Graph) {
	graph.Provide(
		&inject.Object{Value: accountApi.NewClient(os.Getenv("ACCOUNT_ENDPOINT"))},
		&inject.Object{Value: accountClient.NewClient("grpc", os.Getenv("ACCOUNT_ENDPOINT"))},
	)
}

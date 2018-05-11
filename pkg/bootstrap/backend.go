package bootstrap

import (
	"os"

	"github.com/facebookgo/inject"

	"code.ysitd.cloud/api/account"
	"code.ysitd.cloud/api/totp"
	"net/http"
	"time"
)

func injectBackend(graph *inject.Graph) {
	graph.Provide(
		&inject.Object{Value: os.Getenv("ACCOUNT_ENDPOINT"), Name: "account endpoint"},
		&inject.Object{Value: &http.Client{Timeout: 10 * time.Second}},
		&inject.Object{Value: &account.Client{Endpoint: os.Getenv("ACCOUNT_ENDPOINT")}},
		&inject.Object{Value: &totp.Client{Endpoint: os.Getenv("TOTP_ENDPOINT")}},
	)
}

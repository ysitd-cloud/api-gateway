package rest

import (
	"net/http"

	"github.com/gorilla/handlers"

	"golang.ysitd.cloud/log"
	"net/http/httputil"
	"net/url"
)

type Mux struct {
	Frontend http.Handler
	Logger   log.Logger `inject:"http logger"`

	TotpProxy       *TotpProxy `inject:""`
	AccountEndpoint string     `inject:"account endpoint"`
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.Frontend == nil {
		m.initFrontend()
	}

	m.Frontend.ServeHTTP(w, r)
}

func (m *Mux) initFrontend() {
	mux := http.NewServeMux()

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   m.AccountEndpoint,
		Path:   "/",
	})
	mux.Handle("/account/", http.StripPrefix("/account", proxy))
	mux.Handle("/totp/", m.TotpProxy)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unknown route", http.StatusPaymentRequired)
	})

	m.Frontend = handlers.CombinedLoggingHandler(m.Logger.Writer(), mux)
}

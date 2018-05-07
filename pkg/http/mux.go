package http

import (
	"net/http"

	"github.com/gorilla/handlers"

	"golang.ysitd.cloud/log"
)

type Mux struct {
	Frontend http.Handler
	Logger   log.Logger `inject:"http logger"`

	AccountProxy *AccountProxy `inject:""`
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.Frontend == nil {
		m.initFrontend()
	}

	m.Frontend.ServeHTTP(w, r)
}

func (m *Mux) initFrontend() {
	mux := http.NewServeMux()
	mux.Handle("/account/", m.AccountProxy)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unknown route", http.StatusPaymentRequired)
	})

	m.Frontend = handlers.CombinedLoggingHandler(m.Logger.Writer(), mux)
}

package http

import (
	"net/http"
)

type Mux struct {
	Frontend http.Handler

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
	mux.Handle("/account", m.AccountProxy)
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unknown route", http.StatusPaymentRequired)
	}))

	m.Frontend = mux
}

package routes

import (
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/julienschmidt/httprouter"

	"golang.ysitd.cloud/log"

	api "code.ysitd.cloud/api/account"
)

type AccountProxy struct {
	router *httprouter.Router
	Logger log.Logger `inject:"account proxy logger"`
	Backend *api.Client `inject:""`
}

func (p *AccountProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.router == nil {
		p.initRouter()
	}

	p.router.ServeHTTP(w, r)
}

func (p *AccountProxy) initRouter() {
	router := httprouter.New()

	router.POST("/validate", p.validateUserPassword)

	p.router = router
}

func (p *AccountProxy) validateUserPassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when read body", http.StatusBadRequest)
		return
	}

	var input api.ValidateUserRequest

	if err := json.Unmarshal(body, &input); err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when parsing body", http.StatusBadRequest)
		return
	}

	out, err := p.Backend.ValidateUserPassword(r.Context(), &input)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when calling backend", http.StatusBadGateway)
		return
	}

	if !out.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	reply, err := json.Marshal(out.User)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when sending response", http.StatusBadGateway)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(reply)
}

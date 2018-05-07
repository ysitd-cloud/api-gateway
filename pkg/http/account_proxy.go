package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"golang.ysitd.cloud/log"

	api "code.ysitd.cloud/api/account"
	client "code.ysitd.cloud/client/account"
)

type AccountProxy struct {
	router  http.Handler
	Logger  log.Logger         `inject:"account proxy logger"`
	Backend *client.GrpcClient `inject:""`
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
	router.GET("/user/:user", p.getUserInfo)
	router.GET("/token/:token", p.getTokenInfo)

	p.router = http.StripPrefix("/account", router)
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

	out, err := p.Backend.Client.ValidateUserPassword(r.Context(), &input)
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
		http.Error(w, "Error when sending response", 531)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(reply)
}

func (p *AccountProxy) getUserInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	out, err := p.Backend.GetUserInfo(r.Context(), params.ByName("user"))
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when calling backend", http.StatusBadGateway)
		return
	}

	if !out.Exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	reply, err := json.Marshal(out.User)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when sending response", 531)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(reply)
}

func (p *AccountProxy) getTokenInfo(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	out, err := p.Backend.GetTokenInfo(r.Context(), params.ByName("token"))
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when calling backend", http.StatusBadGateway)
		return
	}

	if !out.Exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	reply, err := json.Marshal(out.Token)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when sending response", 531)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(reply)
}

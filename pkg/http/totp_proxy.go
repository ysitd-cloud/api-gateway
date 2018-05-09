package http

import (
	"net/http"

	"golang.ysitd.cloud/log"

	api "code.ysitd.cloud/api/totp"
	"encoding/json"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"time"
)

type TotpProxy struct {
	handler http.Handler
	Logger  log.Logger  `inject:"account proxy logger"`
	Backend *api.Client `inject:""`
}

func (p *TotpProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if p.handler == nil {
		p.initHandler()
	}
}

func (p *TotpProxy) initHandler() {
	router := httprouter.New()

	router.POST("/validate", p.validate)
	router.POST("/key", p.issueKey)
	router.PUT("/key", p.recoverKey)
	router.DELETE("/key", p.removeKey)

	p.handler = http.StripPrefix("/totp", router)
}

func (p *TotpProxy) issueKey(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when read body", http.StatusUnprocessableEntity)
		return
	}

	var input api.IssueKeyRequest

	if err := json.Unmarshal(body, &input); err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when parsing body", http.StatusBadRequest)
		return
	}

	out, err := p.Backend.IssueKey(r.Context(), &input)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when calling backend", http.StatusBadGateway)
		return
	}

	reply, err := json.Marshal(out)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when sending response", 531)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(reply)
}

func (p *TotpProxy) validate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when read body", http.StatusUnprocessableEntity)
		return
	}

	var input struct {
		Issuer   string `json:"issuer"`
		Username string `json:"username"`
		Passcode string `json:"passcode"`
		Time     string `json:"time"`
	}

	if err := json.Unmarshal(body, &input); err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when parsing body", http.StatusBadRequest)
		return
	}

	t, err := time.Parse(time.RFC3339, input.Time)
	if err := json.Unmarshal(body, &input); err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when parsing time", http.StatusUnprocessableEntity)
		return
	}

	out, err := p.Backend.Validate(r.Context(), &api.ValidateRequest{
		Issuer:   input.Issuer,
		Username: input.Username,
		Passcode: input.Passcode,
		Time: &timestamp.Timestamp{
			Seconds: t.Unix(),
			Nanos:   int32(t.Nanosecond()),
		},
	})

	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when calling backend", http.StatusBadGateway)
		return
	}

	if out.Validate {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (p *TotpProxy) recoverKey(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when read body", http.StatusUnprocessableEntity)
		return
	}

	var input api.RecoverRequest

	if err := json.Unmarshal(body, &input); err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when parsing body", http.StatusBadRequest)
		return
	}

	out, err := p.Backend.RecoverKey(r.Context(), &input)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when calling backend", http.StatusBadGateway)
		return
	}

	if out.Validate {
		reply, err := json.Marshal(&api.IssueKeyReply{
			Url:     out.Url,
			Recover: out.Recover,
		})
		if err != nil {
			p.Logger.Error(err)
			http.Error(w, "Error when sending response", 531)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(reply)
	} else {
		w.WriteHeader(http.StatusForbidden)
	}
}

func (p *TotpProxy) removeKey(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when read body", http.StatusUnprocessableEntity)
		return
	}

	var input api.RemoveKeyRequest

	if err := json.Unmarshal(body, &input); err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when parsing body", http.StatusBadRequest)
		return
	}

	out, err := p.Backend.RemoveKey(r.Context(), &input)
	if err != nil {
		p.Logger.Error(err)
		http.Error(w, "Error when calling backend", http.StatusBadGateway)
		return
	}

	if out.Removed {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

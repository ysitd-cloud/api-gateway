package grpc

import (
	"context"

	"bytes"
	api "code.ysitd.cloud/api/account"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func (s *Server) ValidateUserPassword(ctx context.Context, in *api.ValidateUserRequest) (reply *api.ValidateUserReply, err error) {
	body, err := json.Marshal(in)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	req, err := http.NewRequest("POST", s.AccountEndpoint+"/validate", bytes.NewReader(body))
	if err != nil {
		s.Logger.Error(err)
		return
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	reply = new(api.ValidateUserReply)

	err = json.Unmarshal(content, reply)

	if err != nil {
		s.Logger.Error(err)
	}
	return
}

func (s *Server) GetUserInfo(ctx context.Context, in *api.GetUserInfoRequest) (reply *api.GetUserInfoReply, err error) {
	req, err := http.NewRequest("GET", s.AccountEndpoint+"/user/"+in.Username, nil)

	resp, err := s.Client.Do(req)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	reply = new(api.GetUserInfoReply)

	err = json.Unmarshal(content, reply)

	if err != nil {
		s.Logger.Error(err)
	}
	return
}

func (s *Server) GetTokenInfo(ctx context.Context, in *api.GetTokenInfoRequest) (reply *api.GetTokenInfoReply, err error) {
	req, err := http.NewRequest("GET", s.AccountEndpoint+"/token/"+in.Token, nil)

	resp, err := s.Client.Do(req)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.Logger.Error(err)
		return
	}

	reply = new(api.GetTokenInfoReply)

	err = json.Unmarshal(content, reply)

	if err != nil {
		s.Logger.Error(err)
	}
	return
}

package awp

import (
	"awi/config"
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

//hosts 192.168.0.11 avigilon
//https://avigilon:8443/mt/playground
//https://avigilon:8443/mt/playground/rest#!/mt/postLogin

type loginRequest struct {
	Login      string `json:"username"`
	Password   string `json:"password"`
	Token      string `json:"authorizationToken"`
	ClientName string `json:"clientName"` //Не совсем понял - зачем и где это будет использоваться?
}

type loginResponse struct {
	Status string `json:"status"`
	Result result `json:"result"`
}

type result struct {
	Session        string `json:"session"`
	ExternalUserId string `json:"externalUserId"`
	DomainId       string `json:"domainId"`
}

type Auth struct {
	Config   *config.Config
	AuthTime time.Time     // Временная метка последней удачной авторизации
	Request  *loginRequest // Данные для запроса
	Response loginResponse // Ответ от сервера
}

func NewAuth(c *config.Config) *Auth {
	return &Auth{
		Config: c,
		Request: &loginRequest{
			Login:      c.WPUser,
			Password:   c.WPPassword,
			Token:      GenToken(c),
			ClientName: "AWI-Service",
		},
	}
}

func (a *Auth) Login() error {
	//Проверяем, может ещё не стоит авторизоваться снова.
	if time.Since(a.AuthTime).Minutes() < 50 {
		return nil
	}

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(a.Request)
	if err != nil {
		return fmt.Errorf("Login Err: %s", err)
	}
	r := NewRequest(a.Config)
	r.Data = b.Bytes()
	r.Method = POST
	r.Path = "mt/api/rest/v1/login"
	answer, err := r.MakeRequest()
	if err != nil {
		return fmt.Errorf("Login Err: %s", err)
	}

	if err := json.Unmarshal(answer, &a.Response); err != nil {
		return fmt.Errorf("Error decoding config: %s", err)
	}

	if a.Response.Status != "success" {
		return fmt.Errorf("Can't Login: Status == %s", a.Response.Status)
	}
	a.AuthTime = time.Now()

	return nil
}

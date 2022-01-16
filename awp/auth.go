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

type authResponse struct {
	Status string `json:"status"`
	Result struct {
		Session        string `json:"session"`
		ExternalUserId string `json:"externalUserId"`
		DomainId       string `json:"domainId"`
	} `json:"result"`
}

type Auth struct {
	Config   *config.Config
	AuthTime time.Time     // Временная метка последней удачной авторизации
	Request  *loginRequest // Данные для запроса
	Response authResponse  // Ответ от сервера
}

func NewAuth(c *config.Config) *Auth {
	a := &Auth{
		Config: c,
		Request: &loginRequest{
			Login:      c.WPUser,
			Password:   c.WPPassword,
			ClientName: "AWI-Service",
		},
	}

	a.genToken()
	return a
}

func (a *Auth) updateToken() {
	a.Request.Token = a.genToken()
}

//Todo: Добавить мьютекс. Что-бы в случае переподклоючения не потерялся запрос.
func (a *Auth) Login() (*Auth, error) {
	//Проверяем, может ещё не стоит авторизоваться снова.
	timeLimit := 50
	if time.Since(a.AuthTime).Minutes() < float64(timeLimit) {
		return a, nil
	}

	//Дело сессии уже подходит к концу, в первую очередь обновляем токен
	a.updateToken()

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(a.Request)
	if err != nil {
		return a, fmt.Errorf("Login Err: %s", err)
	}
	r := NewRequest(a.Config)
	r.Data = b.Bytes()
	r.Method = POST
	r.Path = "mt/api/rest/v1/login"

	answer, err := r.MakeRequest()
	if err != nil {
		return a, fmt.Errorf("Login: %s", err)
	}

	if err := json.Unmarshal(answer, &a.Response); err != nil {
		return a, fmt.Errorf("Error decoding config: %s", err)
	}

	if a.Response.Status != "success" {
		return a, fmt.Errorf("Can't Login: Status == %s, data: %#v\nAnswer bytes: %s", a.Response.Status, a.Response, string(answer))
	}
	a.AuthTime = time.Now()

	return a, nil
}

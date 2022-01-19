package awp

import (
	"awi/config"
	"sync"
	"time"
)

//hosts 192.168.0.11 avigilon
//https://avigilon:8443/mt/playground
//https://avigilon:8443/mt/playground/rest#!/mt/postLogin
const SessionLimitMinutes float64 = 50

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
	Config        *config.Config
	AuthTime      time.Time     // Временная метка последней удачной авторизации
	LastHeartbeat time.Time     // Когда пришел последний хеартбит
	Request       *loginRequest // Данные для запроса
	Response      authResponse  // Ответ от сервера
	err           error         // Ошибка полученная от WebPointa при авторизации (нет связи, и т.д.)
	wh            *MyWebhooks   // Массив моих ВебХуков
	mu            *sync.Mutex
}

func NewAuth(c *config.Config) *Auth {
	a := &Auth{
		Config: c,
		Request: &loginRequest{
			Login:      c.WPUser,
			Password:   c.WPPassword,
			ClientName: "AWI-Service",
		},
		wh: &MyWebhooks{
			Webhooks: map[string]*Webhook{},
		},
		mu: &sync.Mutex{},
	}

	a.genToken()
	return a
}

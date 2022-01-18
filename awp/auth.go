package awp

import (
	"awi/config"
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
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

func (a *Auth) Lock() {
	a.mu.Lock()
}

func (a *Auth) Unlock() {
	a.mu.Unlock()
}
func (a *Auth) LockState() *sync.Mutex {
	return a.mu
}

// Мьюеткс не трограем
func (a *Auth) updateToken() {
	a.Request.Token = a.genToken()
}

//Todo: Добавить мьютекс. Что-бы в случае переподклоючения не потерялся запрос.
func (a *Auth) Login() (*Auth, error) {
	a.Lock()
	defer a.Unlock()

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
		a.err = fmt.Errorf("Auth.Login Err: %s", err)

		return a, a.err
	}

	r := NewRequest(a.Config)
	r.Data = b.Bytes()
	r.Method = POST
	r.Path = "mt/api/rest/v1/login"

	answer, err := r.MakeRequest()
	if err != nil {
		return a, fmt.Errorf("Auth.Login: %s", err)
	}

	if err := json.Unmarshal(answer, &a.Response); err != nil {
		return a, fmt.Errorf("Auth.Login: Err decoding config: %s", err)
	}

	if a.Response.Status != "success" {
		a.err = fmt.Errorf("Auth.Login: Can't Login: Status == %s, data: %#v\nAnswer bytes: %s", a.Response.Status, a.Response, string(answer))
		return a, a.err
	}

	a.AuthTime = time.Now()

	// Всё хорошо, ошибок нет.
	a.err = nil
	return a, nil
}

func (a *Auth) IsItMyToken(token string) bool {
	a.Lock()
	result := a.wh.IsItMyToken(token)
	a.Unlock()
	return result
}

func (a *Auth) GetError() error {
	//a.Lock()
	err := a.err
	//a.Unlock()
	return err
}

func (a *Auth) UpdateHeartBeat() error {
	a.Lock()
	a.LastHeartbeat = time.Now()
	a.Unlock()
	return nil
}

func (a *Auth) GetHeartBeat() bool {
	//a.Lock()
	var hbState bool
	if time.Since(a.LastHeartbeat).Milliseconds() <= (HeartBeatDelayMs + 100) {
		hbState = true
	}
	//a.Unlock()
	return hbState
}

package awp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func (a *Auth) Login() (*Auth, error) {
	a.Lock()
	defer a.Unlock()

	// Проверяем, может ещё не стоит авторизоваться снова.
	if time.Since(a.AuthTime).Minutes() < SessionLimitMinutes && a.err == nil {
		return a, nil
	}

	// Пробуем закрыть сессию. Даже если нас выбило по-ошибке или потери связи.
	if a.Response.Result.Session != "" {
		if err := a.Logout(a.Response.Result.Session); err != nil {
			a.Response.Result.Session = ""
			a.err = err
			return a, a.err
		}
		log.Printf("Logout: Закрытие сесси [%s] прошло успешно\n", a.Response.Result.Session)
		a.Response.Result.Session = ""
	}

	// Перед логином всегда генерируем новый токен, там присутствует временная метка.
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
	log.Printf("Подключение к WebPoint'у  прошло успешно. SessionId[%s]\n", a.Response.Result.Session)

	a.err = nil
	return a, nil
}

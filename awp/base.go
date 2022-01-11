package awp

import (
	"awi/config"
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Methods string

const (
	GET    Methods = "GET"
	POST   Methods = "POST"
	DELETE Methods = "DELETE"
	PUT    Methods = "PUT"
)

type HttpType string

const (
	HTTP  HttpType = "http"
	HTTPS HttpType = "https"
)

type Request struct {
	Method   Methods
	HType    HttpType          // http|http
	Header   map[string]string // "Content-Type"->'application/json'
	Address  string
	Port     string
	Path     string
	Data     []byte // Данные передаваемые на сервер
	Basic    bool   // Basic авторизация
	Login    string // Логин для BASIC авторизации
	Password string // Пароль для BASIC авторизации
	Response []byte //Ответ от сервера
}

func NewRequest(c *config.Config) *Request {
	return &Request{
		Method:   GET,
		HType:    HTTPS,
		Header:   map[string]string{"Content-Type": "application/json"},
		Address:  c.WPServer,
		Port:     c.WPPort,
		Login:    c.WPUser,
		Password: c.WPPassword,
	}
}

//Makes Request (GET|POST|DELETE|Etc)
func (r *Request) MakeRequest() ([]byte, error) {
	address := r.Address
	if r.Port != "" {
		address += fmt.Sprintf(":%s", r.Port)
	}

	// Генерируем полный путь
	fullPath := fmt.Sprintf("%s://%s/%s", r.HType, address, r.Path)
	req, err := http.NewRequest(string(r.Method), fullPath, bytes.NewReader(r.Data))
	if err != nil {
		return nil, fmt.Errorf("MakeRequest Err: %s", err)
	}

	// Если необходимо - заполняем Header значениями
	for k, v := range r.Header {
		req.Header.Set(k, v)
	}

	// Если нужна basic авторизация
	if r.Basic {
		req.SetBasicAuth(r.Login, r.Password)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("MakeRequest err: %s", err)
	}
	defer resp.Body.Close()

	// Получаем ответ
	answer, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("MakeRequest Err: Response code [%d]%s", resp.StatusCode, resp.Status)
	}
	//fmt.Printf("header: %s\n", resp.Header)
	return answer, nil
}

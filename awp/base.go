package awp

import (
	"awi/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

type ResponseError struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Meta       struct {
		Code int    `json:"code"`
		Name string `json:"name"`
	} `json:"meta"`
}

func ErrorParse(answer []byte) (*ResponseError, error) {
	re := &ResponseError{}
	if err := json.Unmarshal(answer, re); err != nil {
		return nil, fmt.Errorf("ErrorParse: %s", err)
	}
	log.Printf("Response Error: [%d]%s, %s. [code: %d, name: %s]\n", re.StatusCode, re.Status, re.Message, re.Meta.Code, re.Meta.Name)
	return re, nil
}

type Methods string

const (
	GET    Methods = "GET"
	POST   Methods = "POST"
	DELETE Methods = "DELETE"
	PUT    Methods = "PUT"
)

type verbosity string

const (
	LOW    verbosity = "LOW"
	MEDIUM verbosity = "MEDIUM"
	HIGH   verbosity = "HIGH"
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
	c.Lock()
	defer c.Unlock()
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

	//fmt.Printf("rData: %s\n", r.Data)
	//fmt.Printf("fullPath: %s\n", fullPath)

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

	return answer, nil
}

func GenGetter(data map[string]interface{}) string {
	var path string
	var i bool
	for k, v := range data {
		//обработка массива
		if kind, ok := v.(reflect.Kind); ok && kind == reflect.Array {
			for _, arg := range v.([]string) {
				if !i {
					path = fmt.Sprintf("%s=%s", k, arg)
					i = true
				} else {
					fmt.Sprintf("%s&%s=%s", path, k, arg)
				}
			}
			continue
		}
		if !i {
			path = fmt.Sprintf("%s=%s", k, v)
			i = true
		} else {
			path = fmt.Sprintf("%s&%s=%s", path, k, v)
		}
	}
	return path
}

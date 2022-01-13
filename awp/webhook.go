package awp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type webHookPost struct {
	Session string  `json:"session"`
	Webhook Webhook `json:"Webhook"`
}

type ResponseWebhooks struct {
	Status string `json:"status"`
	Result struct {
		Webhooks []Webhook `json:"webhooks"`
		Id       string    `json:"id"`
	} `json:"result"`
}

type RequestWebhooksGet struct {
	Session string   `json:"session"`
	Ids     []string `json:"ids,omitempty"`
}

//Post|Put
type RequestWebhookPost struct {
	Session string   `json:"session"`
	Webhook *Webhook `json:"webhook"`
}

type Webhook struct {
	Id                  string      `json:"id,omitempty"` //Id, полученный после создания
	Url                 string      `json:"url"`
	Heartbeat           Heartbeat   `json:"heartbeat"`
	AuthenticationToken string      `json:"authenticationToken"` // atleast 16 bytes of secure, randomly generated data. Base-64 encoded
	EventTopics         EventTopics `json:"eventTopics"`
}

type myWebhooks struct {
	Webhooks map[string]*Webhook
}

//Используем для мониторигна созданных нами вебхуков
func NewWebhooksMy() myWebhooks {
	return myWebhooks{
		Webhooks: map[string]*Webhook{},
	}
}

type Heartbeat struct {
	Enable      bool `json:"enable"`
	FrequencyMs int  `json:"frequencyMs"`
}

type EventTopics struct {
	WhiteList []string `json:"whitelist,omitempty"`
	BlackList []string `json:"blacklist,omitempty"`
	Include   []string `json:"include,omitempty"`
	Exclude   []string `json:"exclude,omitempty"`
}

func GetWebhooks(a *Auth) ([]Webhook, error) {
	//Всегда проверяем логин перед любым запросом.
	if err := a.Login(); err != nil {
		return nil, fmt.Errorf("GetWebhooks: %s", err)
	}

	query := &RequestWebhooksGet{Session: a.Response.Result.Session}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(query)
	if err != nil {
		return nil, fmt.Errorf("GetWebhooks: %s", err)
	}

	var reqIface map[string]interface{}
	if err := json.NewDecoder(&b).Decode(&reqIface); err != nil {
		return nil, fmt.Errorf("Error decoding reqInface: %s", err)
	}

	r := NewRequest(a.Config)
	r.Data = b.Bytes()
	r.Method = GET
	r.Path = fmt.Sprintf("mt/api/rest/v1/webhooks?%s", GenGetter(reqIface))

	answer, err := r.MakeRequest()
	if err != nil {
		return nil, fmt.Errorf("GetWebhooks: %s", err)
	}

	resp := &ResponseWebhooks{}
	if err := json.Unmarshal(answer, resp); err != nil {
		return nil, fmt.Errorf("Error decoding config: %s", err)
	}

	if resp.Status != "success" {
		d, _ := ErrorParse(answer)
		return nil, fmt.Errorf("Can't read webhooks: Status == %s. [%d]%s - %s", resp.Status, d.StatusCode, d.Status, d.Message)
	}

	return resp.Result.Webhooks, nil
}

func (w *myWebhooks) PostPutWebhook(a *Auth, wh *Webhook, method Methods) error {
	//Всегда проверяем логин перед любым запросом.
	if err := a.Login(); err != nil {
		return fmt.Errorf("PostPutWebhook: %s", err)
	}

	query := &RequestWebhookPost{
		Session: a.Response.Result.Session,
		Webhook: wh,
	}
	//fmt.Printf("Query: %#v\n", query)

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(query)
	if err != nil {
		return fmt.Errorf("PostPutWebhooks: %s", err)
	}
	fmt.Printf("Query encoded: %s\n", b.String())

	r := NewRequest(a.Config)
	r.Data = b.Bytes()
	r.Method = method
	r.Path = "mt/api/rest/v1/webhooks"

	answer, err := r.MakeRequest()
	if err != nil {
		return fmt.Errorf("PostPutWebhooks: %s", err)
	}
	//fmt.Printf("POST WEBHOOK ANSWER: %s\n", string(answer))

	resp := &ResponseWebhooks{}
	if err := json.Unmarshal(answer, resp); err != nil {
		return fmt.Errorf("Error decoding config: %s", err)
	}

	if resp.Status != "success" {
		d, _ := ErrorParse(answer)
		return fmt.Errorf("Can't %s webhooks: Status == %s. [%d]%s - %s", method, resp.Status, d.StatusCode, d.Status, d.Message)
	}

	if method == POST {
		//Заполняем ID, котоый нам выдали при создании
		wh.Id = resp.Result.Id
	}

	//Сохраняем/обновляем вебхук в нашей мапе.
	w.Add(resp.Result.Id, wh)
	return nil
}

func DeleteWebhooks(a *Auth, query *RequestWebhooksGet) error {
	const step = 16
	ids := query.Ids
	l := len(query.Ids)
	for i := 0; ; i += step {
		if i > l {
			break
		}
		border := i + step
		if l < step || border > l {
			border = l
		}

		query.Ids = ids[i:border]
		err := DeleteWebhook(a, query)
		if err != nil {
			return fmt.Errorf("DeleteWebhooks: %s", err)
		}
	}
	return nil
}

func DeleteWebhook(a *Auth, query *RequestWebhooksGet) error {
	//Всегда проверяем логин перед любым запросом.
	if err := a.Login(); err != nil {
		return fmt.Errorf("DeleteWebhook: %s", err)
	}

	query.Session = a.Response.Result.Session
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(query)
	if err != nil {
		return fmt.Errorf("DeleteWebhooks: %s", err)
	}
	//fmt.Printf("Query encoded: %s\n", b.String())

	r := NewRequest(a.Config)
	r.Data = b.Bytes()
	r.Method = DELETE
	r.Path = "mt/api/rest/v1/webhooks"

	answer, err := r.MakeRequest()
	if err != nil {
		return fmt.Errorf("DeleteWebhooks: %s", err)
	}
	//fmt.Printf("DELETE WEBHOOK ANSWER: %s\n", string(answer))

	resp := &ResponseWebhooks{}
	if err := json.Unmarshal(answer, resp); err != nil {
		return fmt.Errorf("Error decoding config: %s", err)
	}

	if resp.Status != "success" {
		d, _ := ErrorParse(answer)
		return fmt.Errorf("Can't DELETE webhooks: Status == %s. [%d]%s - %s", resp.Status, d.StatusCode, d.Status, d.Message)
	}
	return nil
}

func (w *myWebhooks) Add(id string, wh *Webhook) {
	w.Webhooks[id] = wh
}

func (w *myWebhooks) Delete(name string) {
	if _, ok := w.Webhooks[name]; ok {
		delete(w.Webhooks, name)
	}
}

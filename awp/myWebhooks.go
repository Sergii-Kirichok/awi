package awp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type MyWebhooks struct {
	Webhooks map[string]*Webhook
}

//Используем для мониторигна созданных нами вебхуков
//func NewWebhooksMy() *MyWebhooks {
//	return &MyWebhooks{
//		Webhooks: map[string]*Webhook{},
//	}
//}

func (w *MyWebhooks) IsItMyToken(token string) bool {
	for _, webhook := range w.Webhooks {
		if webhook.AuthenticationToken == token {
			return true
		}
	}
	return false
}

func (w *MyWebhooks) Add(id string, wh *Webhook) {
	w.Webhooks[id] = wh
}

func (w *MyWebhooks) Delete(name string) {
	if _, ok := w.Webhooks[name]; ok {
		delete(w.Webhooks, name)
	}
}

func (a *Auth) PostPutWebhook(wh *Webhook, method Methods) error {
	//Всегда проверяем логин перед любым запросом.
	if _, err := a.Login(); err != nil {
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
	//fmt.Printf("Query encoded: %s\n", b.String())

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
		return fmt.Errorf("PostPutWebhooks: err decoding config: %s", err)
	}

	if resp.Status != "success" {
		d, _ := ErrorParse(answer)
		return fmt.Errorf("PostPutWebhooks: Can't %s webhooks: Status == %s. [%d]%s - %s", method, resp.Status, d.StatusCode, d.Status, d.Message)
	}

	if method == POST {
		//Заполняем ID, котоый нам выдали при создании
		wh.Id = resp.Result.Id
	}

	//Сохраняем/обновляем вебхук в нашей мапе.
	a.wh.Add(resp.Result.Id, wh)
	return nil
}

func (wh *MyWebhooks) webhooksReveseCheck(whArr []Webhook) {
	var found bool
	// Перебираем наши вебхуки
	for key := range wh.Webhooks {
		found = false
		for _, v := range whArr {
			// Если нашли в массиве на полученном с сервера - оставляем
			if v.Id == key {
				found = true
			}
		}
		// Если не нашли - удаляем, он какой-то старый
		if !found {
			delete(wh.Webhooks, key)
		}
	}
}

package awp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func (a *Auth) DeleteWebhooks(query *RequestWebhooksGet) error {
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
		err := a.DeleteWebhook(query)
		if err != nil {
			return fmt.Errorf("DeleteWebhooks: %s", err)
		}
	}
	return nil
}

func (a *Auth) DeleteWebhook(query *RequestWebhooksGet) error {
	//Всегда проверяем логин перед любым запросом.
	if _, err := a.Login(); err != nil {
		return fmt.Errorf("DeleteWebhook: %s", err)
	}

	query.Session = a.Response.Result.Session
	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(query)
	if err != nil {
		return fmt.Errorf("DeleteWebhook: %s", err)
	}
	//fmt.Printf("DeleteWebhook: Query encoded: %s\n", b.String())

	r := NewRequest(a.Config)
	r.Data = b.Bytes()
	r.Method = DELETE
	r.Path = "mt/api/rest/v1/webhooks"

	answer, err := r.MakeRequest()
	if err != nil {
		return fmt.Errorf("DeleteWebhook: %s", err)
	}
	//fmt.Printf("DELETE WEBHOOK ANSWER: %s\n", string(answer))

	resp := &ResponseWebhooks{}
	if err := json.Unmarshal(answer, resp); err != nil {
		return fmt.Errorf("DeleteWebhook: err decoding config: %s", err)
	}

	if resp.Status != "success" {
		d, _ := ErrorParse(answer)
		return fmt.Errorf("DeleteWebhook: Can't DELETE webhooks: Status == %s. [%d]%s - %s", resp.Status, d.StatusCode, d.Status, d.Message)
	}
	return nil
}

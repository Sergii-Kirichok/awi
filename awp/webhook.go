package awp

import (
	"awi/config"
	"bytes"
	"encoding/json"
	"fmt"
	"time"
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

func NewWebhook(c *config.Config) *Webhook {
	var addr string
	switch c.WWWPort {
	case "", "443":
		addr = fmt.Sprintf("https://%s/webhooks", c.WWWAddr)
	default:
		addr = fmt.Sprintf("https://%s:%s/webhooks", c.WWWAddr, c.WWWPort)
	}

	return &Webhook{Url: addr,
		Heartbeat: Heartbeat{
			Enable:      true,
			FrequencyMs: 300000, //300000ms -> min value 30sec
		},
		AuthenticationToken: fmt.Sprintf("%x", time.Now()),
		EventTopics: EventTopics{
			WhiteList: []string{
				//"ALL",
				"DEVICE_ANALYTICS_START",
				"DEVICE_ANALYTICS_STOP",
				"DEVICE_DIGITAL_INPUT_ON",
				"DEVICE_DIGITAL_INPUT_OFF",
			},
		},
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
	if _, err := a.Login(); err != nil {
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
		return nil, fmt.Errorf("GetWebhooks: Error decoding reqInface: %s", err)
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
		return nil, fmt.Errorf("GetWebhooks: Error decoding config: %s", err)
	}

	if resp.Status != "success" {
		d, _ := ErrorParse(answer)
		return nil, fmt.Errorf("GetWebhooks: Can't read webhooks: Status == %s. [%d]%s - %s", resp.Status, d.StatusCode, d.Status, d.Message)
	}

	return resp.Result.Webhooks, nil
}

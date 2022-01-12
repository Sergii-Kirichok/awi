package webhooks

import (
	"encoding/json"
	"fmt"
)

func NewMessage() *WebhookMessage {
	return &WebhookMessage{}
}

func (w *WebhookMessage) Parse(b []byte) error {
	if err := json.Unmarshal(b, w); err != nil {
		return fmt.Errorf("Parse: can't unmarshal: %s", err)
	}
	return nil
}

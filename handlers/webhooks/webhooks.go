package webhooks

import (
	"awi/config"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
)

type HandlerData struct {
	cfg *config.Config
	msg *WebhookMessage
}

func NewHandler(cfg *config.Config) *HandlerData {
	return &HandlerData{cfg: cfg, msg: NewMessage()}
}

func (h *HandlerData) WebHooksHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = h.msg.Parse(b)
	if err != nil {
		log.Printf("[ERROR] WebHooksHandler: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.processing()
	if err != nil {
		log.Printf("[ERROR] WebHooksHandler: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	//todo: Remove info messages
	fmt.Printf("Notification parsed [%s], has [%d] events:\n", h.msg.Type, len(h.msg.Notifications))
	for k, e := range h.msg.Notifications {
		fmt.Printf("[INFO] % 2d.Event [%s]=>[%s]:\n\n", k, e.Type, e.Event.Type)
	}
	fmt.Printf("[INFO] ---Data---\n[INFO] %s\n[INFO] ---Data-End---\n\n", string(b))

	w.WriteHeader(http.StatusOK)
}

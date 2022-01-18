package webhooks

import (
	"awi/config"
	"awi/controller"
	_ "embed"
	"io"
	"log"
	"net/http"
)

type HandlerData struct {
	cfg        *config.Config
	controller *controller.Controller
	msg        *WebhookMessage
}

func NewHandler(cfg *config.Config, ctl *controller.Controller) *HandlerData {
	return &HandlerData{
		cfg:        cfg,
		controller: ctl,
		msg:        NewMessage(),
	}
}

func (h *HandlerData) WebHooksHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = h.msg.Parse(b)
	//fmt.Printf("WebHooksHandler: dataIn: %s\n", string(b))
	if err != nil {
		log.Printf("[ERROR] WebHooksHandler: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Обработка входящего вебхука
	err = h.processing()
	if err != nil {
		log.Printf("[ERROR] WebHooksHandler: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//fmt.Printf("[LOG-SART] Webhook Notification parsed [%s], has [%d] events:\n", h.msg.Type, len(h.msg.Notifications))
	//for k, e := range h.msg.Notifications {
	//	fmt.Printf("[INFO] % 2d.Event [%s]=>[%s]:\n", k, e.Type, e.Event.Type)
	//}
	//fmt.Printf("[INFO] ---Data---\n[INFO] %s\n[INFO] ---Data-End---\n", string(b))
	//fmt.Printf("[LOG-SART]\n")

	w.WriteHeader(http.StatusOK)
}

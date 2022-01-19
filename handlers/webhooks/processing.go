package webhooks

import (
	"fmt"
)

// Обработчик ивентов веб-хуков
func (h *HandlerData) processing() error {
	if !h.controller.IsItMyToken(h.msg.AuthenticationToken) {
		return fmt.Errorf("processing: Wrong AuthenticationToken [%s]", h.msg.AuthenticationToken)
	}

	//в зависимости от типа события - выполнить
	switch h.msg.Type {
	case NOTIFICATION:
		for _, e := range h.msg.Notifications {
			if e.Type != "EVENT" { //Обрабатываем только ивенты
				continue
			}
			switch e.Event.Type {
			case DEVICE_DIGITAL_INPUT_ON, DEVICE_DIGITAL_INPUT_OFF:
				return h.inputState(e.Event)
			case DEVICE_ANALYTICS_START:
				return h.personAndCarAnalyticStart(e.Event)
			case DEVICE_ANALYTICS_STOP:
				return h.personAndCarAnalyticStop(e.Event)
			default:
				return fmt.Errorf("processing: NOT Supported evetnt type [%s]", e.Event.Type)
			}
		}
	case HELLO:
		return h.processingHello(h.msg.Type)
	case HEARTBEAT:
		return h.processingHeartbeat(h.msg.Type)
	default:
		return fmt.Errorf("processing: NOT Supported message type [%s]", h.msg.Type)
	}
	return nil
}

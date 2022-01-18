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
			switch EventTypes(e.Event.Type) {
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
		return nil
	case HEARTBEAT:
		return nil
	default:
		return fmt.Errorf("processing: NOT Supported message type [%s]", h.msg.Type)
	}
	return nil
}

//{
//	"siteId":"IN2ir_lQRli_PuW2Un48ZQ",
//	"type":"HELLO",
//	"time":"2022-01-12T13:46:19.981Z",
//	"authenticationToken":"3333746f6b656e3333537472696e67252164284d495353494e4729"
//}
func (h *HandlerData) processingHello() error {
	return nil
}

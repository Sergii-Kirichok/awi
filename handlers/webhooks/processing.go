package webhooks

import (
	"fmt"
)

func (h *HandlerData) processing() error {
	//в зависимости от типа события - выполнить
	for _, e := range h.msg.Notifications {
		switch EventTypes(e.Event.Type) {
		case DEVICE_DIGITAL_INPUT_ON, DEVICE_DIGITAL_INPUT_OFF:
			h.inputState(e.Event)
			return nil
		case DEVICE_ANALYTICS_START, DEVICE_ANALYTICS_STOP:
			h.personAndCarState(e.Event)
			return nil
		default:
			return fmt.Errorf("processing: NOT Supported evetnt type [%s]", e.Event.Type)
		}
	}
	return nil
}

// Устанавливает у входа камеры  (DEVICE_DIGITAL_INPUT) состояние в true|false
func (h *HandlerData) inputState(e *Event) error {
	fmt.Printf("Processing: [%s]\n", e.Type)

	var state bool
	if e.Type == DEVICE_DIGITAL_INPUT_OFF {
		state = true
	}
	errors := map[string]string{}

	for _, cId := range e.CameraIds {
		if err := h.cfg.SetInputState(cId, e.EntityId, state); err != nil {
			errors[cId] = fmt.Sprintf("%s", err)
		}
	}
	if len(errors) > 0 {
		errors := ""
		for cId, err := range errors {
			errors += fmt.Sprintf("inputState: cameraId[%s]:%s\n", cId, err)
		}
		return fmt.Errorf(errors)
	}
	//Если не нашли вход - забиваем, возможно надо-бы ругаться что вход не найден ни у одной отслеживаемой камеры, но нас пока это не волнует
	return nil
}

// Устанавливает у камеры состояние Car и/или Person в true|false
func (h *HandlerData) personAndCarState(e *Event) error {
	fmt.Printf("Processing: [%s]\n", e.Type)
	//Найти камеру по её ID и установить ей состояние входа true
	return nil
}

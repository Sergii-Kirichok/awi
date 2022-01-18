package webhooks

import "fmt"

// Устанавливает у входа камеры  (DEVICE_DIGITAL_INPUT) состояние в true|false
func (h *HandlerData) inputState(e *Event) error {
	fmt.Printf("Processing: [%s]\n", e.Type)

	var state bool
	if e.Type == DEVICE_DIGITAL_INPUT_OFF {
		state = true
	}

	errs := map[string]string{}
	for _, cId := range e.CameraIds {
		if err := h.cfg.SetInputState(cId, e.EntityId, state); err != nil {
			errs[cId] = fmt.Sprintf("%s", err)
		}
	}

	if len(errs) > 0 {
		errrs := ""
		for cId, err := range errs {
			errrs += fmt.Sprintf("inputState: cameraId[%s]:%s\n", cId, err)
		}
		return fmt.Errorf(errrs)
	}
	//Если не нашли вход - забиваем, возможно надо-бы ругаться что вход не найден ни у одной отслеживаемой камеры, но нам пока на это "пи-ли-вать"
	return nil
}

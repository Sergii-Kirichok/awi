package webhooks

import (
	"fmt"
)

// Обработчик ивентов веб-хуков
func (h *HandlerData) processing() error {
	//TODO: Добавить провеку authenticationToken`а
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

// Устанавливает у камеры состояние Car и/или Person в true|false (Обрабатываем только DEVICE_ANALYTICS_START и DEVICE_ANALYTICS_STOP)
// Авижилон должен быть настроен таким образом, что-бы события зоны  по-человеку и по-автомобилю приходили отдельно
func (h *HandlerData) personAndCarAnalyticStart(e *Event) error {
	fmt.Printf("Processing: [%s]\n", e.Type)

	// Обработка событий входа в зону, они имеют массив ClassifiedObjects
	//fmt.Printf("ClassifiedObjects num [%d]\n", len(e.ClassifiedObjects))
	//fmt.Printf("personAndCarAnalyticStart: eventData: %#v\n", e)
	for _, object := range e.ClassifiedObjects {
		switch object.Subclass {
		case VEHICLE, VEHICLE_BICYCLE, VEHICLE_MOTORCYLE, VEHICLE_CAR, VEHICLE_TRUCK, VEHICLE_BUS:
			if e.Activity == OBJECT_PRESENT {
				fmt.Printf("Машина заехала на весовую. - это хорошо, машинка зелёная\n")
				h.cfg.SetCarState(e.CameraId, e.ThisId, true)
				return nil
			}
			return fmt.Errorf("personAndCarAnalyticStart: unsupported vehicle activity: %s", e.Activity)
		case PERSON, PERSON_BODY, PERSON_FACE:
			if e.Activity == OBJECT_PRESENT {
				fmt.Printf("Человек на весовой, это плохо - человечик красный\n")
				h.cfg.SetPersonState(e.CameraId, e.ThisId, false)
				return nil
			}
			return fmt.Errorf("personAndCarAnalyticStart: unsupported person activity: %s", e.Activity)
		}
	}

	// Обработка событий выхода из зоны
	// При выходе из зоны, мы смотрим только на тип ивента (DEVICE_ANALYTICS_STOP), activity (OBJECT_PRESENT) и linkedEventId (carEventId и personEventId)

	return fmt.Errorf("personAndCarAnalyticStart: Wrong event data. Doesn't contain any ClassifiedObjects")
}

func (h *HandlerData) personAndCarAnalyticStop(e *Event) error {
	if e.Activity == OBJECT_PRESENT {
		h.cfg.ClearCarOrPesonState(e.CameraId, e.LinkedEventId)
		return nil
	}
	return fmt.Errorf("personAndCarAnalyticStop: unsupported activity %s", e.Activity)
}

//{
//	"siteId":"IN2ir_lQRli_PuW2Un48ZQ",
//	"type":"HELLO",
//	"time":"2022-01-12T13:46:19.981Z",
//	"authenticationToken":"3333746f6b656e3333537472696e67252164284d495353494e4729"
//}
func (h HandlerData) processingHello() error {
	return nil
}

//{
//	"siteId":"IN2ir_lQRli_PuW2Un48ZQ",
//	"type":"HEARTBEAT",
//	"time":"2022-01-12T16:42:31.349Z",
//	"authenticationToken":"3733746f6b656e3733537472696e67252164284d495353494e4729"
//}
func (h HandlerData) processingHeartbeat() error {
	return nil
}

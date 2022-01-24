package webhooks

import (
	"awi/config"
	"fmt"
)

// Устанавливает у камеры состояние Car и/или Person в true|false (Обрабатываем только DEVICE_ANALYTICS_START и DEVICE_ANALYTICS_STOP)
// Авижилон должен быть настроен таким образом, что-бы события зоны  по-человеку и по-автомобилю приходили отдельно
func (h *HandlerData) personAndCarAnalyticStart(e *Event) error {
	h.ProcessingLogMessage(e)

	// Обработка событий входа в зону, они имеют массив ClassifiedObjects
	//fmt.Printf("ClassifiedObjects num [%d]\n", len(e.ClassifiedObjects))
	//fmt.Printf("personAndCarAnalyticStart: eventData: %#v\n", e)
	for _, object := range e.ClassifiedObjects {
		switch object.Subclass {
		case VEHICLE, VEHICLE_BICYCLE, VEHICLE_MOTORCYLE, VEHICLE_CAR, VEHICLE_TRUCK, VEHICLE_BUS:
			if e.Activity == OBJECT_PRESENT {
				//fmt.Printf("Машина заехала на весовую. - это хорошо, машинка зелёная\n")
				return h.cfg.SetCarState(e.CameraId, e.ThisId, config.StateTrue)
			}
			return fmt.Errorf("personAndCarAnalyticStart: [%s]=>[%s] unsupported activity: [%s]", e.AnalyticEventName, object.Subclass, e.Activity)
		case PERSON, PERSON_BODY, PERSON_FACE:
			if e.Activity == OBJECT_PRESENT {
				//fmt.Printf("Человек на весовой, это плохо - человечик красный\n")
				return h.cfg.SetPersonState(e.CameraId, e.ThisId, config.StateFalse)
			}
			return fmt.Errorf("personAndCarAnalyticStart: [%s]=>[%s] unsupported activity: [%s]", e.AnalyticEventName, object.Subclass, e.Activity)
		}
	}

	// Обработка событий выхода из зоны
	// При выходе из зоны, мы смотрим только на тип ивента (DEVICE_ANALYTICS_STOP), activity (OBJECT_PRESENT) и linkedEventId (carEventId и personEventId)

	return fmt.Errorf("personAndCarAnalyticStart: Wrong event data. Doesn't contain any ClassifiedObjects")
}

func (h *HandlerData) personAndCarAnalyticStop(e *Event) error {
	h.ProcessingLogMessage(e)

	if e.Activity == OBJECT_PRESENT {
		h.cfg.ClearCarOrPesonState(e.CameraId, e.LinkedEventId)
		return nil
	}
	//return fmt.Errorf("personAndCarAnalyticStop: unsupported activity %s", e.Activity)
	return fmt.Errorf("personAndCarAnalyticStart: [%s] unsupported activity: [%s]", e.AnalyticEventName, e.Activity)
}

package webhooks

import (
	"fmt"
	"log"
)

// Устанавливает у камеры состояние Car и/или Person в true|false (Обрабатываем только DEVICE_ANALYTICS_START и DEVICE_ANALYTICS_STOP)
// Авижилон должен быть настроен таким образом, что-бы события зоны  по-человеку и по-автомобилю приходили отдельно
func (h *HandlerData) personAndCarAnalyticStart(e *Event) error {
	log.Printf("Processing: [%s] %s \n", e.Type, e.AnalyticEventName)

	// Обработка событий входа в зону, они имеют массив ClassifiedObjects
	//fmt.Printf("ClassifiedObjects num [%d]\n", len(e.ClassifiedObjects))
	//fmt.Printf("personAndCarAnalyticStart: eventData: %#v\n", e)
	for _, object := range e.ClassifiedObjects {
		switch object.Subclass {
		case VEHICLE, VEHICLE_BICYCLE, VEHICLE_MOTORCYLE, VEHICLE_CAR, VEHICLE_TRUCK, VEHICLE_BUS:
			if e.Activity == OBJECT_PRESENT {
				//fmt.Printf("Машина заехала на весовую. - это хорошо, машинка зелёная\n")
				return h.cfg.SetCarState(e.CameraId, e.ThisId, true)
			}
			return fmt.Errorf("personAndCarAnalyticStart: unsupported vehicle activity: %s", e.Activity)
		case PERSON, PERSON_BODY, PERSON_FACE:
			if e.Activity == OBJECT_PRESENT {
				//fmt.Printf("Человек на весовой, это плохо - человечик красный\n")
				return h.cfg.SetPersonState(e.CameraId, e.ThisId, false)
			}
			return fmt.Errorf("personAndCarAnalyticStart: unsupported person activity: %s", e.Activity)
		}
	}

	// Обработка событий выхода из зоны
	// При выходе из зоны, мы смотрим только на тип ивента (DEVICE_ANALYTICS_STOP), activity (OBJECT_PRESENT) и linkedEventId (carEventId и personEventId)

	return fmt.Errorf("personAndCarAnalyticStart: Wrong event data. Doesn't contain any ClassifiedObjects")
}

func (h *HandlerData) personAndCarAnalyticStop(e *Event) error {
	log.Printf("Processing: [%s] %s \n", e.Type, e.AnalyticEventName)
	if e.Activity == OBJECT_PRESENT {
		h.cfg.ClearCarOrPesonState(e.CameraId, e.LinkedEventId)
		return nil
	}
	return fmt.Errorf("personAndCarAnalyticStop: unsupported activity %s", e.Activity)
}

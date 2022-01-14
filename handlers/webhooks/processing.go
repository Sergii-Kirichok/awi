package webhooks

import "fmt"

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
	//Todo: возмсжно надо боваить отдельный метод для установки состояния входа камеры
	//Todo: ну и наверное всё-же мьютекс необходим
	//Проходим по зонам в поисках камер
	for zId, zone := range h.cfg.Zones {
		// проходим по камерам зоны в поисках подходящей камеры
		for cId, cam := range zone.Cameras {
			//Если нашли камеру, ищем у неё нужный вход
			if cam.Id == e.CameraId {
				//Ищем нужный вход
				for iId, inptut := range cam.Inputs {
					// Нашли вход, дальше будем менять ему статус
					if inptut.EntityId == e.EntityId {
						inptut := h.cfg.Zones[zId].Cameras[cId].Inputs[iId]
						if EventTypes(e.Type) == DEVICE_DIGITAL_INPUT_ON {
							inptut.State = true
							//я-бы вышел через return, но почему-то в сообщении не cameraId, a массив камер...
							continue
						}
						inptut.State = false
					}
				}
			}
		}
	}
	fmt.Printf("[NOTI] ПРОВЕРИТЬ! Не тестировал!\n")
	//Если не нашли вход - забиваем, возможно надо-бы ругаться что вход не найден ни у одной отслеживаемой камеры, но нас пока это не волнует
	return nil
}

// Устанавливает у камеры состояние Car и/или Person в true|false
func (h *HandlerData) personAndCarState(e *Event) error {
	fmt.Printf("Processing: [%s]\n", e.Type)
	//Найти камеру по её ID и установить ей состояние входа true
	return nil
}

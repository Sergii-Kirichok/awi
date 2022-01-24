package config

import (
	"time"
)

// Сбрасываем таймер если не все условия соблюдены
func (c *Config) CountDownZoneCheck(zId string) {
	c.Lock()
	defer c.Unlock()

	for zIndex, zone := range c.Zones {
		if zone.Id == zId {
			resetDuToCam := true

			// Игнорируем статус машины в этой зоне игнора камер
			if zone.IgnoreCarState {
				resetDuToCam = false
			}

			for _, camera := range zone.Cameras {
				// пропускаем пустые камеры
				if camera.Id == "" {
					continue
				}
				// Проверка входов
				for _, input := range camera.Inputs {
					if input.State == StateUnknown || input.State == StateFalse {
						c.Zones[zIndex].TimeLasErr = time.Now()
					}
				}

				// Проверка человека
				if camera.Person == StateUnknown || camera.Person == StateFalse {
					c.Zones[zIndex].TimeLasErr = time.Now()
				}

				// Если статус машинки красный и обычный режим работы - таймер сбрасываем сразу. Или игнорим если в зоне IgnoreCarState
				if (camera.Car == StateUnknown || camera.Car == StateFalse) && !zone.CarOnAnyCamera && !zone.IgnoreCarState {
					c.Zones[zIndex].TimeLasErr = time.Now()
					break
				}

				// Если хоть на одной камере по-машинке всё ок (режим CarOnAnyCamera) - проходим, дальше даже не проверяем
				if (camera.Car == StateUnknown || camera.Car == StateFalse) && zone.CarOnAnyCamera {
					resetDuToCam = false
					break
				}
			}

			// Если режим подтверждения авто хотя-бы на одной камере активен
			if resetDuToCam && zone.CarOnAnyCamera {
				c.Zones[zIndex].TimeLasErr = time.Now()
			}
			// зона с таким Id может быть только одна
			return
		}
	}
}

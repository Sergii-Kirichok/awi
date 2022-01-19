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
			for _, camera := range zone.Cameras {
				// пропускаем пустые камеры
				if camera.Id == "" {
					continue
				}
				// Проверка входов
				for _, input := range camera.Inputs {
					if !input.State {
						c.Zones[zIndex].TimeLasErr = time.Now()
					}
				}

				// Проверка человека
				if !camera.Person {
					c.Zones[zIndex].TimeLasErr = time.Now()
				}

				// Если статус машинки красный и обычный режим работы - таймер сбрасываем сразу
				if !camera.Car && !zone.CarOnAnyCamera {
					c.Zones[zIndex].TimeLasErr = time.Now()
					break
				}
				// Если хоть на одной камере по-машинке всё ок - проходим, дальше даже не проверяем
				if camera.Car && zone.CarOnAnyCamera {
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

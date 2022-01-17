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
			for _, camera := range zone.Cameras {
				// пропускаем пустые камеры
				if camera.Id == "" {
					continue
				}
				for _, input := range camera.Inputs {
					if !input.State {
						c.Zones[zIndex].TimeLasErr = time.Now()
					}
				}
				// Машина должна быть, а человек должен отсутствовать
				if !camera.Car || !camera.Person {
					c.Zones[zIndex].TimeLasErr = time.Now()
				}
			}
		}
	}
}

package config

import (
	"fmt"
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
						fmt.Printf("Воход генерирует тревожноге событие и сбрасывает таймер\n")
						c.Zones[zIndex].TimeLasErr = time.Now()
					}
				}
				// Машина должна быть, а человек должен отсутствовать
				if !camera.Car || !camera.Person {
					fmt.Printf("[%d]%s Camera: %s => Person: [%s]%v, Car: [%s]%v\n", zIndex, zone.Name, camera.Name, camera.PersonEventId, camera.Person, camera.CarEventId, camera.Car)
					c.Zones[zIndex].TimeLasErr = time.Now()
				}
			}
		}
	}
}

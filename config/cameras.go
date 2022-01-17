package config

import "fmt"

// Если будем использовать не только при старте, добавить мьютексы и проверить на deadlock
func (c *Config) camerasCheckInTheZone(zoneIndex int) (bool, error) {
	var needUpdate bool
	for index, camera := range c.Zones[zoneIndex].Cameras {
		if camera.Serial == "" {
			needUpdate = true
			//удаляем камеры с пустым серийным номером
			if len(c.Zones[zoneIndex].Cameras) == 1 {
				c.Zones[zoneIndex].Cameras = []Cam{}
				continue
			}
			c.Zones[zoneIndex].Cameras[index] = c.Zones[zoneIndex].Cameras[len(c.Zones[zoneIndex].Cameras)-1]
			c.Zones[zoneIndex].Cameras[len(c.Zones[zoneIndex].Cameras)-1] = Cam{}
			c.Zones[zoneIndex].Cameras = c.Zones[zoneIndex].Cameras[:len(c.Zones[zoneIndex].Cameras)-1]
		}
	}
	return needUpdate, nil
}

// Устанавливаем значение указанного входа  для  указанной камеры
func (c *Config) SetInputState(cid string, inId string, state bool) error {
	c.Lock()
	defer c.Unlock()
	for zId, zone := range c.Zones {
		for camId, camera := range zone.Cameras {
			if camera.Id == cid {
				for index, input := range camera.Inputs {
					if input.EntityId == inId {
						c.Zones[zId].Cameras[camId].Inputs[index].State = state
						return nil
					}
				}
			}
		}
	}
	return fmt.Errorf("SetInputState: Camera.Id[%s] doesn't has input [%s]\n", cid, inId)
}

// Устанавливаем статус для  камеры по eё Id
func (c *Config) SetCarState(cid string, eventId string, state bool) error {
	c.Lock()
	defer c.Unlock()
	for zId, zone := range c.Zones {
		for camId, camera := range zone.Cameras {
			if camera.Id == cid {
				c.Zones[zId].Cameras[camId].Car = state
				c.Zones[zId].Cameras[camId].CarEventId = eventId
				return nil
			}
		}
	}
	return fmt.Errorf("SetCarState: Can't finde Camera.Id[%s] in any zone\n", cid)
}

// Устанавливаем статус для  человека по id камеры
func (c *Config) SetPersonState(cid string, eventId string, state bool) error {
	c.Lock()
	defer c.Unlock()
	for zId, zone := range c.Zones {
		for camId, camera := range zone.Cameras {
			if camera.Id == cid {
				c.Zones[zId].Cameras[camId].Person = state
				c.Zones[zId].Cameras[camId].PersonEventId = eventId
				return nil
			}
		}
	}
	return fmt.Errorf("SetCarState: Can't finde Camera.Id[%s] in any zone\n", cid)
}

// Вызывается по-окончанию события (DEVICE_ANALYTICS_STOP).
// Для машины это false - Когда машина выехала из зоны, она станет красной
// Для человека это true - когда человек вышел из зоны это хорошо и человечик должен стать зелёным
func (c *Config) ClearCarOrPesonState(cid string, eventId string) error {
	c.Lock()
	defer c.Unlock()
	for zId, zone := range c.Zones {
		for camId, camera := range zone.Cameras {
			if camera.Id == cid {
				if c.Zones[zId].Cameras[camId].CarEventId == eventId {
					c.Zones[zId].Cameras[camId].CarEventId = ""
					c.Zones[zId].Cameras[camId].Car = false
					return nil
				}
				if c.Zones[zId].Cameras[camId].PersonEventId == eventId {
					c.Zones[zId].Cameras[camId].PersonEventId = ""
					c.Zones[zId].Cameras[camId].Person = true
					return nil
				}

			}
		}
	}
	return fmt.Errorf("SetCarState: Can't finde Camera.Id[%s] and/or eventId [%s] in any zone\n", cid, eventId)
}

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

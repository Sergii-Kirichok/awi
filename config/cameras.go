package config

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

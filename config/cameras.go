package config

func (c *Config) camerasCheckInTheZone(zoneId int) (bool, error) {
	var needUpdate bool
	for index, camera := range c.Zones[zoneId].Cameras {
		if camera.Serial == "" {
			needUpdate = true
			//удаляем камеру с пустым серийным номером
			//c.Zones[zoneId].Cameras = append(c.Zones[zoneId].Cameras[:index], c.Zones[zoneId].Cameras[index+1:]...)
			if len(c.Zones[zoneId].Cameras) == 1 {
				c.Zones[zoneId].Cameras = []cam{}
				continue
			}
			c.Zones[zoneId].Cameras[index] = c.Zones[zoneId].Cameras[len(c.Zones[zoneId].Cameras)-1]
			c.Zones[zoneId].Cameras[len(c.Zones[zoneId].Cameras)-1] = cam{}
			c.Zones[zoneId].Cameras = c.Zones[zoneId].Cameras[:len(c.Zones[zoneId].Cameras)-1]

		}
	}
	return needUpdate, nil
}

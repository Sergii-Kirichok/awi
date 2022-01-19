package controller

// Отсюда веб берёт данные по Zone, со всеми её статусами
func (c *Controller) GetZoneData(zoneId string) Zone {
	c.mu.Lock()
	defer c.mu.Unlock()

	for zId, zData := range c.zones {
		if zId == zoneId {
			//log.Printf("controller.getZoneData: Name:%s, heartBeat: %v, Error: %s\n", zData.Name, zData.Heartbeat, zData.Error)
			return *zData
		}
	}
	return Zone{Error: "зона с таким ID відсутня"}
}

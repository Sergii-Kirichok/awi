package config

import "fmt"

const zoneMinDelaySec int = 30

func (c *Config) checkZones() error {
	var willUpdate bool
	for zId, z := range c.Zones {
		// Проверка ID
		zName := genZoneId(z.Name)
		if z.Id != zName {
			willUpdate = true
			c.Zones[zId].Id = zName
		}
		// Проверка минимального времени задержки
		if zoneMinDelaySec > z.DelaySec {
			willUpdate = true
			c.Zones[zId].DelaySec = zoneMinDelaySec
		}
		// Проверка метода сохранения.
		if !z.Bookmarks && !z.Alarms {
			willUpdate = true
			c.Zones[zId].Bookmarks = true
		}
		// Проверка камер в зоне
		camsUpdated, err := c.camerasCheckInTheZone(zId)
		if err != nil {
			return fmt.Errorf("checkZones: %s", err)
		}
		if camsUpdated {
			willUpdate = true
		}
	}
	//ID Зон обновлены, сохраняем конфиг
	if willUpdate {
		if err := c.Save(); err != nil {
			return fmt.Errorf("checkZones: %s", err)
		}
	}
	return nil
}

func genZoneId(name string) string {
	return fmt.Sprintf("%x%x", ZoneNameAppendix, name)
}

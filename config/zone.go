package config

import (
	"crypto/sha256"
	"fmt"
)

const zoneMinDelaySec int = 30

// Если будем использовать не только при старте - добавить использование мьютексов
func (c *Config) checkZones() error {
	var willUpdate bool
	for zIndex, z := range c.Zones {
		// Проверка ID
		zName := genZoneId(z.Name)
		if z.Id != zName {
			willUpdate = true
			c.Zones[zIndex].Id = zName
		}
		// Проверка минимального времени задержки
		if zoneMinDelaySec > z.DelaySec {
			willUpdate = true
			c.Zones[zIndex].DelaySec = zoneMinDelaySec
		}
		// Проверка метода сохранения.
		if !z.Bookmarks && !z.Alarms {
			willUpdate = true
			c.Zones[zIndex].Bookmarks = true
		}
		// Проверка камер в зоне
		camsUpdated, err := c.camerasCheckInTheZone(zIndex)
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
	data := fmt.Sprintf("%s%s", ZoneNameAppendix, name)
	hash := sha256.Sum256([]byte(data))
	return fmt.Sprintf("%x", hash)
}

// Дергаем в для получения копии текущих даннызх зоны
func (c *Config) GetZoneData(zId string) Zone {
	c.Lock()
	defer c.Unlock()
	for _, z := range c.Zones {
		if z.Id == zId {
			return z
		}
	}
	return Zone{}
}

func (c *Config) GetZoneNames() map[string]string {
	c.Lock()
	names := make(map[string]string, len(c.Zones))
	for _, z := range c.Zones {
		names[z.Id] = z.Name
	}
	c.Unlock()
	return names
}

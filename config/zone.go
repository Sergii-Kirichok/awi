package config

import (
	"crypto/sha256"
	"fmt"
	"time"
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

// Дергаем в для получения копии текущих данных зоны
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

func ZoneIsOk(z *Zone) bool {
	if time.Since(z.TimeLasErr).Seconds() < float64(z.DelaySec) {
		return false
	}

	var carState bool
	// Режим когда не важно есть машинка или нет
	if z.IgnoreCarState {
		carState = true
	}

	for _, cam := range z.Cameras {
		for _, input := range cam.Inputs {
			if !input.State || !cam.Person {
				return false
			}
		}
		// Если не в режиме машина на любой камере и на камере нет машины
		if !z.CarOnAnyCamera && !cam.Car {
			return false
		}
		// Если есть машинка хоть на одной камере ставим ок
		if z.CarOnAnyCamera && cam.Car {
			carState = true
		}
	}

	// Нет ни одной машины
	if !carState {
		return false
	}

	return true
}

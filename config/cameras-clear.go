package config

import (
	"fmt"
	"time"
)

// не лочим мьютекс, т.к. вызывающая этот метод функция уже залочила конфиг
func (c *Config) ClearCamStates(Id string) error {
	for zId, zone := range c.Zones {
		for camId, camera := range zone.Cameras {
			if camera.Id == Id {
				c.Zones[zId].Cameras[camId].Car = ""
				c.Zones[zId].Cameras[camId].CarEventId = ""
				c.Zones[zId].Cameras[camId].Person = ""
				c.Zones[zId].Cameras[camId].PersonEventId = ""
				for inpIndex := range camera.Inputs {
					c.Zones[zId].Cameras[camId].Inputs[inpIndex].State = ""
				}
				c.Zones[zId].TimeLasErr = time.Now()
				return nil
			}
		}
	}

	return fmt.Errorf("ClearCamStates: Can't find cam with this ID [%s]", Id)
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
					c.Zones[zId].Cameras[camId].Car = StateFalse
					c.Zones[zId].TimeLasErr = time.Now()
					return nil
				}
				if c.Zones[zId].Cameras[camId].PersonEventId == eventId {
					c.Zones[zId].Cameras[camId].PersonEventId = ""
					c.Zones[zId].Cameras[camId].Person = StateTrue
					c.Zones[zId].TimeLasErr = time.Now()
					return nil
				}
			}
		}
	}
	return fmt.Errorf("SetCarState: Can't finde Camera.Id[%s] and/or eventId [%s] in any zone\n", cid, eventId)
}

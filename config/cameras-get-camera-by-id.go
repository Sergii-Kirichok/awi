package config

import "fmt"

func (c *Config) GetCameraById(camId string) (Cam, error) {
	c.Lock()
	defer c.Unlock()
	for _, zone := range c.Zones {
		for _, camera := range zone.Cameras {
			if camera.Id == camId {
				return camera, nil
			}
		}
	}
	return Cam{}, fmt.Errorf("Can't find cam by id [%s]", camId)
}

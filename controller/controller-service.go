package controller

import "time"

// Рутина собирающая данные из конфига, обрабатывающая и генерирующая данные для отображения веба.
// Веб, получает данные по-зонам обращаясь к методам Controller`a
func (c *Controller) Service() {
	for {
		c.auth.Lock()
		confNames := c.auth.Config.GetZoneNames()
		for zId := range confNames {
			c.auth.Config.CountDownZoneCheck(zId)
			c.updateZone(zId) //можно сделать рутинкой.
		}
		c.auth.Unlock()
		time.Sleep(1 * time.Second)
	}
}

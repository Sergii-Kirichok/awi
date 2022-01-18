package controller

func (c *Controller) UpdateHeartBeat() error {
	c.mu.Lock()
	err := c.auth.UpdateHeartBeat()
	c.mu.Unlock()
	return err
}

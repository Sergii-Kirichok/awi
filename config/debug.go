package config

func (c *Config) GetDebug() bool {
	c.Lock()
	debug := c.Debug
	c.Unlock()
	return debug
}

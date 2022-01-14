package config

func (c *Config) Lock() {
	c.mu.Lock()
}

func (c *Config) Unlock() {
	c.mu.Unlock()
}

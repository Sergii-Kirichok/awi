package config

//Set default dev data
func (c *Config) checkNonce() {
	c.Lock()
	c.Unlock()
	if c.DevNonce == "" || c.Key == "" {
		c.Nonce = "FO#09121901"
		c.Key = "9fbd5669d18031f8ce5d4261b17dc3334c78f9e1597bef0bb5d3c26c7cffee8a"
		return
	}
	c.Nonce = c.DevNonce
	c.Key = c.DevKey
}

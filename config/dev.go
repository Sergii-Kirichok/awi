package config

//Set default dev data
func (c *Config) checkNonce() {
	if c.DevNonce == "" || c.DevKey == "" {
		c.DevNonce = "FO#09121901"
		c.DevKey = "9fbd5669d18031f8ce5d4261b17dc3334c78f9e1597bef0bb5d3c26c7cffee8a"
	}
}

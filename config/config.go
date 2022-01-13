package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func New() *Config {
	return &Config{}
}

//Load configuration data from the encrypted file
func (c *Config) Load(dir string) {
	file := fmt.Sprintf("%s/config.json", dir)
	if _, err := os.Stat(file); err == nil {
		data, _ := ioutil.ReadFile(file)
		if err := json.Unmarshal(data, &c); err != nil {
			log.Fatalf("Error decoding config: %s", err)
		}
	} else if os.IsNotExist(err) {
		c.makeDefault()
		c.Save(dir, "config.json")
	} else {
		log.Fatalf("Error reading config: %s", err)
	}
}

//Encrypt configuration and Save it
func (c *Config) Save(dir string, fileName string) {
	f, err := os.OpenFile(fmt.Sprintf("%s/%s", dir, fileName), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	defer f.Close()

	if err != nil {
		log.Fatalf("Error writing config: %s", err)
	}

	out := json.NewEncoder(f)
	out.SetIndent("", "\t")

	if err := out.Encode(c); err != nil {
		log.Fatalf("Error encoding config: %s", err)
	}
}

func (c *Config) makeDefault() {
	c.WWWAddr = "127.0.0.1"
	c.WWWPort = "8080"
	c.WWWCertificate = "certificates/src_certs_server.pem"
	c.WWWCertificateKey = "certificates/src_certs_server.key"
	c.WPServer = "avigilon"
	c.WPPort = "8443"
	c.WPUser = "administrator"
	c.WPPassword = "yjdsqgfhjkm"
	c.DevNonce = "FO#09121901"
	c.DevKey = "9fbd5669d18031f8ce5d4261b17dc3334c78f9e1597bef0bb5d3c26c7cffee8a"
	c.Zones = []zone{
		0: {
			Name:      "Весовая №1",
			Bookmarks: true,
			Alarms:    true,
			DelaySec:  300,
			Cameras: []cam{
				0: {
					Serial: "102109218992",
					Name:   "Фронтальная 1.1",
				},
				1: {
					Serial: "1234567890123",
					Name:   "Фронтальная 1.2",
				},
			},
		},
		1: {
			Name:      "Весовая №2",
			Bookmarks: true,
			Cameras: []cam{
				0: {
					Serial: "1234567890124",
					Name:   "Фронтальная 2.1",
				},
				1: {
					Serial: "1234567890125",
					Name:   "Фронтальная 2.2",
				},
			},
		},
	}
}

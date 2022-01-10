package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

//Generate default confiig, encrypt it and store to the default file (config.bin)
func makeDefaultConfig(fileName string) []byte {
	defConf := []byte(`{
		"www_addr":"",
		"www_http_https_redirect":true,
		"www_certificate":"certificates/src_certs_server.pem",
		"www_certificate_key":"certificates/src_certs_server.key",
		"www_session_time_life":"30",
		"db_server":"127.0.0.1",
		"db_port":"1433",
		"db_name":"Visitors",
		"db_user":"userDefault",
		"db_password":"123",
		"rest_server":"127.0.0.1:9797",
		"rest_user":"admin",
		"rest_password":"ResT_AdmiN",
		"rest_sync_interval":"60",
		"rest_outside_location":"L1",
		"rest_outside_area":"A1"
}`)

	if err := ioutil.WriteFile(fileName, defConf, 0644); err != nil {
		log.Printf("Error writing default config: %s", err)
	}
	return defConf
}

//Load configuration data from the encrypted file
func LoadConfiguration(dir string) Config {
	var config Config
	file := fmt.Sprintf("%s/config.json", dir)
	if _, err := os.Stat(file); err == nil {
		data, _ := ioutil.ReadFile(file)
		if err := json.Unmarshal(data, &config); err != nil {
			log.Fatalf("Error decoding config: %s", err)
		}
	} else if os.IsNotExist(err) {
		if err := json.Unmarshal(makeDefaultConfig(file), &config); err != nil {
			log.Fatalf("Error parsing default config: %s", err)
		}
	} else {
		log.Fatalf("Error reading config: %s", err)
	}
	return config
}

//Encrypt configuration and Save it
func SaveConfiguration(fileName string, cfg Config) {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("Error writing config: %s", err)
	}
	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		log.Fatalf("Error encoding config: %s", err)
	}
	if err := f.Close(); err != nil {
		log.Printf("Error closing config: %s", err)
	}
}

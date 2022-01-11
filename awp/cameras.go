package awp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type Camera struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Available       bool   `json:"available"`
	ConnectionState string `json:"connectionState"`
	FirmwareVersion string `json:"firmwareVersion"`
	IpAddress       string `json:"ipAddress"`
	Model           string `json:"model"`
	Serial          string `json:"serial"`
	ServerId        string `json:"serverId"`
	PhysicalAddress string `json:"physicalAddress"`
	Active          bool   `json:"active"`
	LogicalId       int    `json:"logicalId"`
	links           []struct {
		Type   string `json:"type"` //"DIGITAL_INPUT","DIGITAL_OUTPUT","AUDIO_INPUT","AUDIO_OUTPUT",
		Id     string `json:"id"`
		Source string `json:"source"`
		Target string `json:"target"`
	} `json:"links"`
	Capabilities     map[string][]string `json:"capabilities"`
	Connected        bool                `json:"connected"`
	ConnectionStatus struct {
		IsConnectable bool          `json:"isConnectable"`
		State         string        `json:"state"`
		ErrorFlags    []interface{} `json:"errorFlags"`
		StartTime     time.Time     `json:"startTime"`
	} `json:"connectionStatus"`
}

type ResponseCameras struct {
	Status string `json:"status"`
	Result struct {
		Cameras []Camera `json:"cameras"`
	} `json:"result"`
}

type RequestCameras struct {
	Session   string    `json:"session"`
	Verbosity verbosity `json:"verbosity"`
}

// Возвращает список доступных камер
func GetCameras(a *Auth) ([]Camera, error) {
	//Всегда проверяем логин перед любым запросом.
	if err := a.Login(); err != nil {
		return nil, fmt.Errorf("GetCameras: %s", err)
	}

	query := &RequestCameras{
		Session:   a.Response.Result.Session,
		Verbosity: HIGH,
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(query)
	if err != nil {
		return nil, fmt.Errorf("GetCameras: %s", err)
	}

	var reqIface map[string]interface{}
	if err := json.NewDecoder(&b).Decode(&reqIface); err != nil {
		return nil, fmt.Errorf("Error decoding reqInface: %s", err)
	}

	r := NewRequest(a.Config)
	r.Data = b.Bytes()
	r.Method = GET
	r.Path = fmt.Sprintf("mt/api/rest/v1/cameras?%s", GenGetter(reqIface))

	answer, err := r.MakeRequest()
	if err != nil {
		return nil, fmt.Errorf("GetCameras: %s", err)
	}

	resp := &ResponseCameras{}
	if err := json.Unmarshal(answer, resp); err != nil {
		return nil, fmt.Errorf("Error decoding config: %s", err)
	}

	if resp.Status != "success" {
		d, _ := ErrorParse(answer)
		return nil, fmt.Errorf("Can't read cameras: Status == %s. [%d]%s - %s", resp.Status, d.StatusCode, d.Status, d.Message)
	}

	return resp.Result.Cameras, nil
}

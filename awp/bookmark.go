package awp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type RequestBookmark struct {
	Session  string    `json:"session"`
	Bookmark BookmarkT `json:"bookmark"`
}

type ResponseBookmark struct {
	Status string `json:"status"`
	Result struct {
		Id string `json:"id"`
	} `json:"result"`
}

type BookmarkT struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CameraIds   []string  `json:"cameraIds"`
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	IsProtected bool      `json:"isProtected"`
}

func (a *Auth) MakeBookmark(ZoneId string) (*ResponseBookmark, error) {

	//Всегда проверяем логин перед любым запросом.
	if _, err := a.Login(); err != nil {
		return nil, fmt.Errorf("GetCameras: %s", err)
	}

	zone := a.Config.GetZoneData(ZoneId)
	if zone.Name == "" {
		return nil, fmt.Errorf("MakeBookmark: Can't make bookmark. Unknown zoneId %s", ZoneId)
	}

	query := &RequestBookmark{
		Session: a.Response.Result.Session,
		Bookmark: BookmarkT{
			Name:        fmt.Sprintf("%s - Зважування", zone.Name),
			StartTime:   time.Now().Add(-10 * time.Second),
			EndTime:     time.Now(),
			CameraIds:   []string{},
			IsProtected: true,
		},
	}

	for _, camera := range zone.Cameras {
		if camera.Id != "" {
			query.Bookmark.CameraIds = append(query.Bookmark.CameraIds, camera.Id)
		}
	}

	var b bytes.Buffer
	err := json.NewEncoder(&b).Encode(query)
	if err != nil {
		return nil, fmt.Errorf("MakeBookmark: %s", err)
	}

	//var reqIface map[string]interface{}
	//if err := json.NewDecoder(&b).Decode(&reqIface); err != nil {
	//	return nil, fmt.Errorf("MakeBookmark: err decoding reqInface: %s", err)
	//}

	r := NewRequest(a.Config)
	r.Data = b.Bytes()
	r.Method = POST
	//r.Path = fmt.Sprintf("mt/api/rest/v1/bookmark?%s", GenGetter(reqIface))
	r.Path = "mt/api/rest/v1/bookmark"

	answer, err := r.MakeRequest()
	if err != nil {
		return nil, fmt.Errorf("MakeBookmark: %s", err)
	}

	//fmt.Printf("MakeBookmark: Answer: %s\n", string(answer))
	resp := &ResponseBookmark{}
	if err := json.Unmarshal(answer, resp); err != nil {
		return nil, fmt.Errorf("MakeBookmark: err decoding config: %s", err)
	}

	if resp.Status != "success" {
		d, _ := ErrorParse(answer)
		return nil, fmt.Errorf("MakeBookmark: Can't read cameras: Status == %s. [%d]%s - %s", resp.Status, d.StatusCode, d.Status, d.Message)
	}
	return resp, nil
}

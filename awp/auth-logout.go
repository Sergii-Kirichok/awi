package awp

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type requestLogOut struct {
	Session string `json:"session"`
}

func (a *Auth) Logout(ssesionId string) error {
	b := new(bytes.Buffer)

	query := &requestLogOut{
		Session: a.Response.Result.Session,
	}

	err := json.NewEncoder(b).Encode(query)
	if err != nil {
		return fmt.Errorf("Auth.Logout Err: %s", err)
	}

	r := NewRequest(a.Config)
	r.Data = b.Bytes()
	r.Method = POST
	r.Path = "mt/api/rest/v1/logout"

	answer, err := r.MakeRequest()
	if err != nil {
		return fmt.Errorf("Auth.Logout: %s", err)
	}

	if err := json.Unmarshal(answer, &a.Response); err != nil {
		return fmt.Errorf("Auth.Logout: Err decoding config: %s", err)
	}

	if a.Response.Status != "success" {
		return fmt.Errorf("Auth.Logout: Can't Logout: Status == %s, data: %#v\nAnswer bytes: %s", a.Response.Status, a.Response, string(answer))
	}

	return nil
}

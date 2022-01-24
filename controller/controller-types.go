package controller

import (
	"awi/awp"
	"awi/config"
	"sync"
)

type Controller struct {
	auth  *awp.Auth
	zones map[string]*Zone
	mu    *sync.Mutex
}

type Zone struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Heartbeat   bool               `json:"heartbeat"`
	TimeLeftSec int                `json:"timeLeft"`
	Cameras     map[string]*Camera `json:"cameras"`
	Error       string             `json:"error"`
}

type Camera struct {
	Id         string            `json:"id"`
	Name       string            `json:"name"`
	Human      config.States     `json:"human,omitempty"`
	Car        config.States     `json:"car,omitempty"`
	Inputs     map[string]*Input `json:"inputs"`
	Connection Connection        `json:"connection"`
}

type Input struct {
	Id    string        `json:"id"`
	State config.States `json:"state,omitempty"` // red,green и  "" == gray
}

type Connection struct {
	Type  config.CamState `json:"type"`  // CONNECTED, etc...
	State bool            `json:"state"` // true == green, false == red
}

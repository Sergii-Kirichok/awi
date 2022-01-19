package controller

import (
	"awi/awp"
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
	Id     string            `json:"id"`
	Name   string            `json:"name"`
	Human  bool              `json:"human"`
	Car    bool              `json:"car"`
	Inputs map[string]*Input `json:"inputs"`
}

type Input struct {
	Id    string `json:"id"`
	State bool   `json:"state"`
}

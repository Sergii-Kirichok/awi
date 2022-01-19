package controller

import (
	"awi/awp"
	"sync"
)

func New(a *awp.Auth) *Controller {
	a.Config.Lock()
	zonesNum := len(a.Config.Zones)
	a.Config.Unlock()
	return &Controller{
		auth:  a,
		mu:    &sync.Mutex{},
		zones: make(map[string]*Zone, zonesNum),
	}
}

func (c *Controller) IsItMyToken(token string) bool {
	return c.auth.IsItMyToken(token)
}

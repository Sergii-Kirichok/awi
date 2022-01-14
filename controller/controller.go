package controller

import (
	"awi/config"
	"sync"
)

type Controller struct {
	conf  *config.Config
	zones map[string]*Zone
	mu    *sync.Mutex
}

type Zone struct {
	Name        string
	TimeLeftSec int
	Cameras     map[string]Camera
}

type Camera struct {
	Id     string
	Name   string
	Human  bool
	Car    bool
	Inputs map[string]Input
}

type Input struct {
	Id    string
	State bool
}

func New(c *config.Config) *Controller {
	return &Controller{
		conf:  c,
		mu:    &sync.Mutex{},
		zones: make(map[string]*Zone, len(c.Zones)),
	}
}

// Рутина собирающая данные из конфига, обрабатывающая и генерирующая данные для отображения веба.
// Веб, получает данные по-зонам обращаясь к методам Controller`a
func (c *Controller) Service() {
	for {

	}
}

func (c *Controller) updateZone(name string) {
	zConf := c.conf.GetZoneData(name)
	z := c.zones[zConf.Name]
	z.Name = zConf.Name
	z.TimeLeftSec = zConf.TimeLeft.Second()
}

// Веб берёт данные по Zone, со всеми её статусами
func (c *Controller) GetZoneData(name string) Zone {
	c.mu.Lock()
	zone := *c.zones[name]
	c.mu.Unlock()
	return zone
}

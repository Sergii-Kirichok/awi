package controller

import (
	"awi/config"
	"sync"
	"time"
)

type Controller struct {
	conf  *config.Config
	zones map[string]*Zone
	mu    *sync.Mutex
}

type Zone struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	TimeLeftSec int                `json:"time_left_sec"`
	Cameras     map[string]*Camera `json:"cameras"`
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
		confNames := c.conf.GetZoneNames()
		for zId := range confNames {
			c.conf.CountDownZoneCheck(zId)
			c.updateZone(zId)
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *Controller) updateZone(zId string) {
	zConf := c.conf.GetZoneData(zId)
	c.mu.Lock()
	z, ok := c.zones[zConf.Id]
	c.mu.Unlock()

	// Если первый запуски и зоны такой нет - создаём
	if !ok {
		z = &Zone{
			Id:   zConf.Id,
			Name: zConf.Name,
		}
	}

	// Обновляем счётчик времени
	timeSince := int(time.Since(zConf.TimeLasErr).Seconds())
	if timeSince >= zConf.DelaySec {
		z.TimeLeftSec = 0
	} else {
		z.TimeLeftSec = zConf.DelaySec - timeSince
	}

	// Проверка есть-ли вообще у зоны мапа камер
	if len(z.Cameras) == 0 {
		z.Cameras = map[string]*Camera{}
	}

	// Обновляем данные по-камерам
	for _, cam := range zConf.Cameras {
		// Данные по камере отсутствуют в системе
		if cam.Id == "" {
			continue
		}
		if _, ok := z.Cameras[cam.Id]; !ok {
			z.Cameras[cam.Id] = &Camera{}
		}
		z.Cameras[cam.Id].Name = cam.Name
		z.Cameras[cam.Id].Car = cam.Car
		z.Cameras[cam.Id].Human = cam.Person
		z.Cameras[cam.Id].Id = cam.Id

		// Проверяем, есть-ли вообще у камеры мапа входов
		if z.Cameras[cam.Id].Inputs == nil {
			//Проверяем, может на камере вообще нет входов (т.е. когда мы получали данные в sync-eре webpoint нам сказал что у камеры нет входов нужного типа)
			z.Cameras[cam.Id].Inputs = make(map[string]*Input)
		}

		// Заполняем данные о статусе входов (если таковые имеются)
		for _, cInput := range cam.Inputs {
			if _, ok := z.Cameras[cam.Id].Inputs[cInput.EntityId]; !ok {
				z.Cameras[cam.Id].Inputs[cInput.EntityId] = &Input{}
			}
			z.Cameras[cam.Id].Inputs[cInput.EntityId].Id = cInput.EntityId
			z.Cameras[cam.Id].Inputs[cInput.EntityId].State = cInput.State
		}
	}

	c.mu.Lock()
	c.zones[zConf.Id] = z
	c.mu.Unlock()
}

// Отсюда веб берёт данные по Zone, со всеми её статусами
func (c *Controller) GetZoneData(zoneId string) Zone {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Поиск в реально существующей зоны, если зоны нет - отдадим пустую
	for zId, zData := range c.zones {
		if zId == zoneId {
			return *zData
		}
	}
	return Zone{TimeLeftSec: 36000}
}

func (c *Controller) MakeAction(zoneId string) error {
	return nil
}

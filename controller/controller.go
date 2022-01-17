package controller

import (
	"awi/awp"
	"errors"
	"sync"
	"time"
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
	Webpoint    bool               `json:"webpoint"`
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

// Рутина собирающая данные из конфига, обрабатывающая и генерирующая данные для отображения веба.
// Веб, получает данные по-зонам обращаясь к методам Controller`a
func (c *Controller) Service() {
	for {
		c.auth.Lock()

		confNames := c.auth.Config.GetZoneNames()
		for zId := range confNames {
			c.auth.Config.CountDownZoneCheck(zId)
			c.updateZone(zId) //можно сделать рутинкой.
		}
		c.auth.Unlock()
		time.Sleep(1 * time.Second)
	}
}

func (c *Controller) updateZone(zId string) {
	zConf := c.auth.Config.GetZoneData(zId)
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

var zoneErr = errors.New("zone doesn't exist")

// Отсюда веб берёт данные по Zone, со всеми её статусами
func (c *Controller) GetZoneData(zoneId string) (Zone, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//todo: Добавить сюда-же вывод других ошибок. например потеря связи с WP, проблема с авторизацией и т.д...
	// Поиск в реально существующей зоны, если зоны нет - отдадим пустую
	for zId, zData := range c.zones {
		if zId == zoneId {
			return *zData, nil
		}
	}

	return Zone{}, zoneErr
}

func (c *Controller) MakeAction(zoneId string) error {
	return nil
}

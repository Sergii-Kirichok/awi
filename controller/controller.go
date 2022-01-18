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
	TimeLeftSec int                `json:"timeLeft"`
	Cameras     map[string]*Camera `json:"cameras"`
	Err         error              `json:"error"` // todo: mb error string
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

	// Отдаём WebPopint connection status
	if err := c.auth.GetError(); err != nil {
		z.Webpoint = false
		z.Err = err
	} else {
		// Отдаём heartBeat status только если с авторизацией всё хорошо ...
		z.Heartbeat = c.auth.GetHeartBeat()
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
	for zId, zData := range c.zones {
		if zId == zoneId {
			return *zData, zData.Err
		}
	}

	return Zone{}, zoneErr
}

func (c *Controller) MakeAction(zoneId string) error {
	//fmt.Printf("ASK FOR MAKE ACTION FOR ZONE: %s\n", zoneId)
	_, err := c.auth.MakeBookmark(zoneId)
	return err
}

//func (c Config) MakeAction(zoneId string) error {
//	z := c.GetZoneData(zoneId)
//	if z.Alarms {
//		fmt.Printf("Config.MakeAction: Sending alarm to WebPoint Zone: %s\n", zoneId)
//		//return err
//	}
//	if z.Bookmarks {
//		fmt.Printf("Config.MakeAction: Sending Bookmarks to WebPoint. Zone: %s\n", zoneId)
//		//return err
//	}
//	return nil
//}

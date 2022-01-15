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
	Name        string
	TimeLeftSec int
	Cameras     map[string]*Camera
}

type Camera struct {
	Id     string
	Name   string
	Human  bool
	Car    bool
	Inputs map[string]*Input
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
		confNames := c.conf.GetZoneNames()
		for _, name := range confNames {
			c.updateZone(name)
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *Controller) updateZone(name string) {
	zConf := c.conf.GetZoneData(name)

	z, ok := c.zones[zConf.Name]
	// Если первый запуски и зоны такой нет - создаём
	if !ok {
		z = &Zone{}
	}

	z.Name = zConf.Name
	z.TimeLeftSec = zConf.TimeLeft.Second()

	// Проверка есть-ли вообще у зоны мапа камер
	if z.Cameras == nil {
		l := len(zConf.Cameras)
		// Проверяем, есть-ли у этой зоны вообще камеры. Может администратор создал зону в конфише, а камеры в неё ещё не добавил...
		if l > 0 {
			z.Cameras = make(map[string]*Camera, len(zConf.Cameras))
		}
	}

	// Обновляем данные по-камерам
	for _, cam := range zConf.Cameras {
		zCam := z.Cameras[cam.Id]
		zCam.Name = cam.Name
		zCam.Car = cam.Car
		zCam.Human = cam.Person
		zCam.Id = cam.Id

		// Проверяем, есть-ли вообще у камеры мапа входов
		if zCam.Inputs == nil {
			//Проверяем, может на камере вообще нет входов (т.е. когда мы получали данные в sync-eре webpoint нам сказал что у камеры нет входов нужного типа)
			l := len(cam.Inputs)
			if l > 0 {
				zCam.Inputs = make(map[string]*Input)
			}
		}

		// Заполняем данные о статусе входов (если таковые имеются)
		for _, cInput := range cam.Inputs {
			zCam.Inputs[cInput.EntityId].Id = cInput.EntityId
			zCam.Inputs[cInput.EntityId].State = cInput.State
		}
	}
}

// Веб берёт данные по Zone, со всеми её статусами
func (c *Controller) GetZoneData(name string) Zone {
	c.mu.Lock()
	// Поиск в реально существующей зоны, если зоны нет - отдадим пустую
	for zName, zData := range c.zones {
		if zName == name {
			return *zData
		}
	}
	c.mu.Unlock()
	return Zone{}
}

func (c *Controller) MakeAction(name string) error {

	return nil
}

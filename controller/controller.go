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
	Id          string
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
			//fmt.Printf("Zone: [%s]%s\n", zId, name)
			c.updateZone(name)
		}
		time.Sleep(1 * time.Second)
	}
}

func (c *Controller) updateZone(zId string) {
	zConf := c.conf.GetZoneData(zId)

	z, ok := c.zones[zConf.Id]
	// Если первый запуски и зоны такой нет - создаём
	if !ok {
		z = &Zone{
			Id:   zConf.Id,
			Name: zConf.Name,
		}
	}

	z.TimeLeftSec = zConf.TimeLeft.Second()

	//TODO: Разобраться с логикой камер
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
		z.Cameras[cam.Id].Name = cam.Name
		z.Cameras[cam.Id].Car = cam.Car
		z.Cameras[cam.Id].Human = cam.Person
		z.Cameras[cam.Id].Id = cam.Id

		// Проверяем, есть-ли вообще у камеры мапа входов
		if z.Cameras[cam.Id].Inputs == nil {
			//Проверяем, может на камере вообще нет входов (т.е. когда мы получали данные в sync-eре webpoint нам сказал что у камеры нет входов нужного типа)
			l := len(cam.Inputs)
			if l > 0 {
				z.Cameras[cam.Id].Inputs = make(map[string]*Input)
			}
		}

		// Заполняем данные о статусе входов (если таковые имеются)
		for _, cInput := range cam.Inputs {
			z.Cameras[cam.Id].Inputs[cInput.EntityId].Id = cInput.EntityId
			z.Cameras[cam.Id].Inputs[cInput.EntityId].State = cInput.State
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
	return Zone{TimeLeftSec: 36000}
}

func (c *Controller) MakeAction(name string) error {

	return nil
}

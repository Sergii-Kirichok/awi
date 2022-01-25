package controller

import (
	"awi/config"
	"time"
)

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

		// Данные о статусе камеры и её каким цветом она будет отображена (зелёный/красный) в вебе
		var state bool
		if cam.ConState == config.CamConnected {
			state = true
		}
		z.Cameras[cam.Id].Connection = Connection{
			Type:  cam.ConState,
			State: state,
		}

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

	// Формируем ошибку если надо.
	z.Error = ""
	if err := c.auth.GetError(); err != nil {
		z.Error = err.Error()
	}

	// Отдаём heartBeat status только если с авторизацией всё хорошо ...
	z.Heartbeat = c.auth.GetHeartBeat()

	c.mu.Lock()
	c.zones[zConf.Id] = z
	c.mu.Unlock()
}

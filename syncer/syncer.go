package syncer

import (
	"awi/awp"
	"fmt"
	"log"
	"sync"
	"time"
)

type syncer struct {
	auth     *awp.Auth
	wh       *awp.MyWebhooks // Массив моих вебхуков. По-хорошему вебхук вообще должен быть один
	lastSinc time.Time       // Время последней синхронизации
	blocker  bool            // Блокировка работы основной программы, если нет возможности синхронизироваться ok == true
	m        *sync.Mutex
}

func New(a *awp.Auth) *syncer {
	return &syncer{
		auth: a,
		wh:   awp.NewWebhooksMy(),
		m:    &sync.Mutex{},
	}
}

// Основная рутина
func (s *syncer) Sync() {
	var showInfo bool
	for {
		// Заполняем главную структуру актуальными данными
		if time.Since(s.lastSinc).Seconds() > 60 {
			s.lastSinc = time.Now()
			//Обновляем/заполняем данными камеры и входы
			if err := s.update(); err != nil {
				log.Printf("[ERROR] Can't sync data: %s\n", err)
			}

			// Удаляем старые/чужие вебхуки если таковые имеются, создаём свои если их ешё нет
			if s.getBlocker() {
				if err := awp.WebhooksUpdater(s.auth, s.wh); err != nil {
					log.Printf("[ERROR] %s\n", err)
				}
			}
		}

		if !showInfo {
			showInfo = true
			fmt.Printf("MyWebHooks: %v, data: %#v\n", s.wh, s.wh.Webhooks)
			for k, v := range s.wh.Webhooks {
				fmt.Printf("Webhooks Key %v, Value: %#v\n", k, v)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}

// Обвновляем данные между веб-поинтом (сервером авижилона) и нашей основной структурой (данные о зонах, камерах, входах)
// Вебхуки мы обрабатываем на стороне веб-сервера
// если блокер == false, значит нет связи с WebPoint'om
func (s *syncer) update() error {
	s.m.Lock()
	defer s.m.Unlock()
	// Забираем у WebPointa все камеры
	cameras, err := awp.GetCameras(s.auth)
	if err != nil {
		s.blocker = true
		return fmt.Errorf("update: %s", err)
	}
	//todo: удалить инфо-вывод
	for _, cam := range cameras {
		log.Printf("[INFO] Camera Name: \"%s\", Serial: %s, Active: %v, Id: %s\n", cam.Name, cam.Serial, cam.Active, cam.Id)
	}
	//todo: Обновляем камеры (имя,id и т.д.) в нашем конфиге (рабочей структуре)

	//todo: Обновляем входы на камерах, прописываем их ИД

	s.blocker = true // Всё ок, блокер должен отключен == ok == true
	return nil
}

func (s *syncer) getBlocker() bool {
	return s.blocker
}

//fmt.Printf("MyWebHooks: %v, data: %#v\n", wh, wh.Webhooks)
//	for k, v := range wh.Webhooks {
//		fmt.Printf("Webhooks Key %v, Value: %#v\n", k, v)
//		v.Heartbeat.FrequencyMs = 30000
//
//		if err := wh.PostPutWebhook(auth, v, awp.PUT); err != nil {
//			fmt.Printf("POST Error: %s\n", err)
//		}
//	}

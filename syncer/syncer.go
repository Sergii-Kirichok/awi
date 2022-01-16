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
	lastSync time.Time       // Время последней синхронизации
	blocker  blocker         // Блокировка работы основной программы, если нет возможности синхронизироваться ok == true
	m        *sync.Mutex
}
type blocker bool

const (
	Blocked           blocker = false
	NotBlocked        blocker = true
	WebPointSyncDelay float64 = 10
)

func New(a *awp.Auth) *syncer {
	return &syncer{
		auth: a,
		wh:   awp.NewWebhooksMy(),
		m:    &sync.Mutex{},
	}
}

// Основная рутина синхронизирующая данные между веб-понитом и нами
func (s *syncer) Sync() {
	var showInfo bool
	for {
		// Заполняем главную структуру актуальными данными
		if time.Since(s.lastSync).Seconds() > WebPointSyncDelay {
			s.m.Lock()
			s.lastSync = time.Now()
			s.m.Unlock()
			// Обновляем/заполняем данными камеры и входы
			if err := s.update(); err != nil {
				log.Printf("[ERROR] Sync: Can't sync data: %s\n", err)
			}

			// Удаляем старые/чужие вебхуки если таковые имеются, создаём свои если их ешё нет
			if s.GetBlocker() == NotBlocked {
				if err := awp.WebhooksUpdater(s.auth, s.wh); err != nil {
					log.Printf("[ERROR] Sync: %s\n", err)
				}
			}

		}

		if !showInfo && s.auth.Config.GetDebug() {
			showInfo = true
			fmt.Printf("[INFO] MyWebHooks: %v, data: %#v\n", s.wh, s.wh.Webhooks)
			for k, v := range s.wh.Webhooks {
				fmt.Printf("[INFO] Webhooks Key %v, Value: %#v\n", k, v)
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

	// Забираем у WebPointa все доступные камеры
	cameras, err := awp.GetCameras(s.auth)
	if err != nil {
		s.blocker = Blocked
		return fmt.Errorf("update: %s", err)
	}

	//todo: удалить инфо-вывод
	// Обновляем данные камеры и ёё входов (имя,id и т.д.) в нашей рабочей структуре
	for _, camera := range cameras {
		//log.Printf("[INFO] Camera Name: \"%s\", Serial: %s, Active: %v, Id: %s\n", camera.Name, camera.Serial, camera.Active, camera.Id)
		//log.Printf("[INFO] Camera FULL DATA: %#v\n", camera)
		//log.Printf("[INFO] Config: %#v\n", s.auth.Config)

		s.cameraSet(camera)
	}
	s.blocker = NotBlocked // Всё ок
	return nil
}

func (s *syncer) GetBlocker() blocker {
	s.m.Lock()
	bState := s.blocker
	s.m.Unlock()
	return bState
}

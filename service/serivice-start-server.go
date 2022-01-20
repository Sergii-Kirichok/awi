package service

import (
	"awi/awp"
	"awi/config"
	"awi/controller"
	"awi/syncer"
	"awi/webserver"
	"fmt"
	"time"
)

func startServer(s *Myservice) {
	s.Info(1, fmt.Sprintf("Запуск сервера %s", s))

start:
	cfg, err := config.New().Load()
	if err != nil {
		s.Error(1, fmt.Sprintf("Load config: %s", err))
		time.Sleep(10 * time.Second)
		goto start
	}

	//Если не смогли авторизоваться при старте, то явно какая-то ошибка в конфиге или с самим сервером. Ругаемся, но идём дальше (нужен веб для отображения информации)
	auth, err := awp.NewAuth(cfg).Login()
	if err != nil {
		s.Error(2, fmt.Sprintf("%s", err))
	}

	// Сохраняем auth, что-бы дернуть LogOut по закрытию программы
	s.Auth = auth

	// Синхронизатор. Проверяет конфиг, находит Ид камер, создаёт и удаляет вебхуки, ...
	go syncer.New(auth).Sync()

	// Cтруктура хранящая данные используемые для генерации данных по зонам, камерам и ошибок (связи,синхронизации) передаваемых при get-запросе из веба
	control := controller.New(auth)

	// Сервис отвечающий за: таймеры обратного отсчёта по зонам и их состояния. Веб работает с ней и методами controllera
	go control.Service()

	// WebServer, принимает и обрабатываем webhook-и от WebPointa, так-же отдаёт страничку с Кнопкой, таймером обратного отсчёта, значками состояния, ...
	webserver.New(s.Name, s.Version, cfg, control).ListenAndServeHTTPS()
}

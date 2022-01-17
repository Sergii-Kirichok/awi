package main

import (
	"awi/awp"
	"awi/config"
	"awi/controller"
	"awi/syncer"
	"awi/webserver"
	"log"
	"time"
)

func main() {
start:
	cfg, err := config.New().Load()
	if err != nil {
		log.Printf("[ERROR] Load config: %s\n", err)
		time.Sleep(10 * time.Second)
		goto start
	}

	//Если не смогли авторизоваться при старте, то явно какая-то ошибка в конфиге или с самим сервером WP. Поэтому идём перечитываем конфиг и пробуем ещё раз.
	auth, err := awp.NewAuth(cfg).Login()
	if err != nil {
		log.Printf("[ERROR] Auth.Login: %s", err)
	}

	// Синхронизатор. Проверяет конфиг, находит Ид камер, создаёт и удаляет вебхуки, ...
	go syncer.New(auth).Sync()

	//структура хранящая данные для генерации данных по зонам, камерам и ошибок (связи,синхронизации) передаваемых при get-запросе из веба
	control := controller.New(auth)
	//сервис отвечающий за: таймеры обратного отсчёта по зонам и их состояния. Веб работает с ней и  методами controllera\
	go control.Service()

	// WebServer, принимает и обрабатываем webhook-и от WebPointa, так-же отдаёт страничку с Кнопкой, таймером обратного отсчёта, значками состояния, ...
	webserver.New("Avigilon Weight Integration Server", "beta 0.1", cfg, control).ListenAndServeHTTPS()

}

package main

import (
	"awi/awp"
	"awi/config"
	"awi/syncer"
	"awi/webserver"
	"log"
)

func main() {
	cfg := config.New().Load()

	auth := awp.NewAuth(cfg)
	if err := auth.Login(); err != nil {
		log.Fatalf("Login Err: %s", err)
	}

	// Синхронизатор. Проверяет конфиг, находит Ид камер, содаёт и удаляет вебхуки, ...
	s := syncer.New(auth)
	go s.Sync()

	//Todo: Тут рутинка, отвечающая за, таймеры обратного отсчёта по зонам

	// WebServer, принимает и обрабатываем webhook-и от WebPointa, так-же отдаёт страничку с Кнопкой, таймером обратного отсчёта, значками состояния, ...
	web := webserver.New("Avigilon Weight Integration Server", "beta 0.1", cfg)
	web.ListenAndServeHTTPS()
}

package main

import (
	"awi/awp"
	"awi/config"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

func main() {
	rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	cfg := config.New()
	cfg.Load(rootDir)

	auth := awp.NewAuth(cfg)
	if err := auth.Login(); err != nil {
		log.Fatalf("Login Err: %s", err)
	}
	//fmt.Printf("Auth Data: %#v", auth)

	// Получаем все камеры
	cameras, err := awp.GetCameras(auth)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	for _, cam := range cameras {
		fmt.Printf("Camera Name: \"%s\", Serial: %s, Active: %v\n", cam.Name, cam.Serial, cam.Active)
	}

	// Получаем все активные вебхуки на сервере
	webhooks, err := awp.GetWebhooks(auth)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	if len(webhooks) > 0 {
		var ids []string
		for _, hook := range webhooks {
			fmt.Printf("[Will be deleted] WebHook ID[%s]: URL: \"%s\", HeartBeat: %v, Events: %v\n", hook.Id, hook.Url, hook.Heartbeat, hook.EventTopics)
			ids = append(ids, hook.Id)
		}
		awp.DeleteWebhook(auth, &awp.RequestWebhooksGet{Ids: ids})
	}

	//Создаём вебхук
	rand.Seed(time.Now().UnixNano())
	rnd := rand.Intn(100)
	wh := awp.NewWebhooksMy()
	webhook := &awp.Webhook{
		Url: "https://192.168.0.10/webhooks",
		Heartbeat: awp.Heartbeat{
			Enable:      true,
			FrequencyMs: 300000, //300000ms = 5 min
		},
		AuthenticationToken: fmt.Sprintf("%x", fmt.Sprintf("%dtoken%dString%d", rnd, rnd)),
		EventTopics: awp.EventTopics{
			WhiteList: []string{"ALL"},
		},
	}
	wh.PostPutWebhook(auth, webhook, awp.POST)

	fmt.Printf("MyWebHooks: %v, data: %#v\n", wh, wh.Webhooks)
	for k, v := range wh.Webhooks {
		fmt.Printf("Webgooks Key %v, Value: %#v\n", k, v)
		v.Heartbeat.FrequencyMs = 100000

		if err := wh.PostPutWebhook(auth, v, awp.PUT); err != nil {
			fmt.Printf("POST Error: %s\n", err)
		}
	}
}

package awp

import (
	"fmt"
	"log"
)

// Удаляем лишние вебхуки. Оставляем только те, которые есть у нас в всписке
func WebhooksUpdater(auth *Auth, wh *MyWebhooks) error {
	// Получаем все активные вебхуки на сервере
	webhooks, err := GetWebhooks(auth)
	if err != nil {
		return fmt.Errorf("WebhookClear: %s", err)
	}

	// Если нашли что-то - будем проверять
	if len(webhooks) > 0 {
		ids := webhooksWPCheck(wh, webhooks)
		if len(ids) > 0 {
			if err := DeleteWebhooks(auth, &RequestWebhooksGet{Ids: ids}); err != nil {
				return fmt.Errorf("WebhooksUpdater: %s\n", err)
			}
		}

	}

	// Обратная проверке, есть-ли в нашем массиве лишние вебхуки (отсутствующие). Еcли есть - там-же их и удаляем.
	if len(wh.Webhooks) > 0 {
		//todo: Добавить проверку на соответствие параметров вебхука полученного от WP нашим настройкам (может кто-то сделал PUT и их изменил)
		wh.webhooksReveseCheck(webhooks)
	} else {
		// нет на сервере вебхуков, надо создать
		webhook := NewWebhook(auth.Config)
		wh.PostPutWebhook(auth, webhook, POST)
	}
	//Провека, все-ли вебхуки мы создали, если не хватает - добавить. Хотя, вообще-то, должен быть всего один!
	return nil
}

// Проверка вебхуков полученных от WebPointa
func webhooksWPCheck(wh *MyWebhooks, whArr []Webhook) []string {
	//log.Printf("[INFO] Found %d webhooks on wp\n", len(whArr))
	var ids []string
	for _, hook := range whArr {
		if _, ok := wh.Webhooks[hook.Id]; !ok {
			log.Printf("[INFO] Going to delete this WebHook -> ID[%s]: URL: \"%s\", HeartBeat: %v, Events: %#v\n", hook.Id, hook.Url, hook.Heartbeat, hook.EventTopics)
			ids = append(ids, hook.Id)
		}
	}
	return ids
}

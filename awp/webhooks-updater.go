package awp

import (
	"fmt"
	"log"
)

const webhooksTolerance = 4

// Удаляем лишние вебхуки. Оставляем только те, которые есть у нас в всписке
func (a *Auth) WebhooksUpdater() error {
	// Получаем все активные вебхуки на сервере WebPointа (внутри функция логин, a.Lock вызовет deadlock)
	webhooks, err := a.GetWebhooksFromWP()
	if err != nil {
		return fmt.Errorf("WebhooksUpdater: %s", err)
	}

	// Если нашли что-то - будем проверять
	if len(webhooks) > webhooksTolerance {
		ids := a.webhooksWPCheck(webhooks)
		if len(ids) > {
			if err := a.DeleteWebhooks(&RequestWebhooksGet{Ids: ids}); err != nil {
				return fmt.Errorf("WebhooksUpdater: %s\n", err)
			}
		}
	}

	a.Lock()
	whLen := len(a.wh.Webhooks)
	a.Unlock()

	// Обратная проверке, есть-ли в нашем массиве лишние вебхуки (отсутствующие на сервере). Еcли есть - там-же их и удаляем.
	if whLen > 0 {
		// Todo: Добавить проверку на соответствие параметров вебхука полученного от WP нашим настройкам (может кто-то сделал PUT и их изменил)
		a.Lock()
		a.wh.webhooksReveseCheck(webhooks)
		a.Unlock()
	} else {
		// нет на сервере вебхуков, надо создать
		webhook := NewWebhook(a.Config)
		if err := a.PostPutWebhook(webhook, POST); err != nil {
			return fmt.Errorf("WebhooksUpdater: %s\n", err)
		}
	}
	// Проверка, все-ли вебхуки мы создали, если не хватает - добавить. Хотя, вообще-то, должен быть всего один!
	return nil
}

// Проверка вебхуков полученных от WebPointa. На выходе массив лишних вебхуков
func (a *Auth) webhooksWPCheck(whArr []Webhook) []string {
	a.Lock()
	defer a.Unlock()

	// Сравниваем полученные вебхуки с нашим массивом
	var ids []string
	for _, hook := range whArr {
		if _, ok := a.wh.Webhooks[hook.Id]; !ok {
			log.Printf("[INFO] Going to delete this WebHook -> ID[%s]: URL: \"%s\", HeartBeat: %v, Events: %#v\n", hook.Id, hook.Url, hook.Heartbeat, hook.EventTopics)
			ids = append(ids, hook.Id)
		}
	}
	return ids
}

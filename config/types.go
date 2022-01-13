package config

import (
	"sync"
	"time"
)

type Config struct {
	WWWAddr           string      `json:"www_addr"`            // Адрес на котором принимаем запросы
	WWWPort           string      `json:"www_Port"`            // Порт на котором принимаем запросы
	WWWCertificate    string      `json:"www_certificate"`     // Путь к файлу сертификатов https
	WWWCertificateKey string      `json:"www_certificate_key"` // Путь к файлу ключа сертификата
	WPServer          string      `json:"wp_server"`           // webPoint  адрес
	WPPort            string      `json:"wp_port"`             // webPoint  порт
	WPUser            string      `json:"wp_user"`             // webPoint  user
	WPPassword        string      `json:"wp_password"`         // webPoint  password
	DevNonce          string      `json:"dev_nonce,omitempty"` // Avigilon developer nonce
	DevKey            string      `json:"dev_key,omitempty"`   // Avigilon developer key
	Zones             []zone      `json:"zones"`               // Список весовых зон
	Debug             bool        `json:"debug,omitempty"`     // Работа в режиме отладки
	mu                *sync.Mutex `json:"-"`
}

type zone struct {
	Name      string    `json:"name"`                // Имя зоны -> в вебе будет использоваться для отображения (?zone=base64(name))
	Cameras   []cam     `json:"cameras"`             // Камеры в пределах текущей зоны
	DelaySec  int       `json:"delay_sec"`           // Задержка после сработки входа, наличия машины и отсутствия человека
	Bookmarks bool      `json:"bookmarks,omitempty"` // Генерировать Закладки
	Alarms    bool      `json:"alarms,omitempty"`    // Генерировать тревоги
	State     bool      `json:"state,omitempty"`     // Текущее состояние (красная/зелёная) (результирующий - человек, машина, вход, задержка)
	TimeLeft  time.Time `json:"-"`                   // Время, которое осталось до активации кнопки взвешивания
	Countdown bool      `json:"-"`                   // Можно-ли начинать обратный отсчёт по зоне.
}

type cam struct {
	CamID    string  `json:"-"`                // ИД-Камеры. Получаем по RESTу на основании serial
	Serial   string  `json:"serial"`           // Серийный номер камеры, по нему и ёё и идентифицируем
	ConState string  `json:"connectionState"`  // 'CONNECTED'
	Name     string  `json:"_"`                // Имя камеры
	Inputs   []input `json:"-"`                // Состояние входов
	Car      bool    `json:"car,omitempty"`    // В зоне обнаружена машина
	Person   bool    `json:"person,omitempty"` // В зоне обнаружен человек
}

type input struct {
	EntityId string `json:"-"` // Заполняем динамически, берём у камеры в links []{ {type:"DIGITAL_INPUT", source: "4xIx1DMwMLSwMDW2tDBKNNBLycwzMBASCDilIfJR0W3apqrIovO_tncAAA"},{},...}
	State    bool   `json:"-"` // Статус входа
}

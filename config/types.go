package config

import (
	"sync"
	"time"
)

const ZoneNameAppendix string = "ZoneNameIs"

type Config struct {
	WWWAddr           string      `json:"www_addr"`            // Адрес на котором принимаем запросы
	WWWPort           string      `json:"www_Port"`            // Порт на котором принимаем запросы
	WWWCertificate    string      `json:"www_certificate"`     // Путь к файлу сертификатов https
	WWWCertificateKey string      `json:"www_certificate_key"` // Путь к файлу ключа сертификата
	WPServer          string      `json:"wp_server"`           // webPoint  адрес
	WPPort            string      `json:"wp_port"`             // webPoint  порт
	WPUser            string      `json:"wp_user"`             // webPoint  user
	WPPassword        string      `json:"wp_password"`         // webPoint  password
	Nonce             string      `json:"-"`                   // Avigilon developer nonce
	DevNonce          string      `json:"dev_nonce,omitempty"` // DevNonce if we need it
	Key               string      `json:"-"`                   // Avigilon developer key
	DevKey            string      `json:"dev_key,omitempty"`   // DevKey if we need it
	Zones             []Zone      `json:"zones"`               // Список весовых зон
	Debug             bool        `json:"debug,omitempty"`     // Работа в режиме отладки
	mu                *sync.Mutex `json:"-"`
}

type Zone struct {
	Id         string    `json:"Id"`              // По нему будем работать с Зоной.?Zone=hexEncoded(ZoneNameAppendix+name) (Обязательно обновлять и сохранять в конфиге если при чтении конфига была пустая)
	Name       string    `json:"name"`            // Имя зоны -> в вебе будет использоваться для отображения (?Zone=hexEncoded(ZoneNameAppendix+name))
	Cameras    []Cam     `json:"cameras"`         // Камеры в пределах текущей зоны
	DelaySec   int       `json:"delay_sec"`       // Задержка после сработки входа, наличия машины и отсутствия человека
	State      bool      `json:"state,omitempty"` // Текущее состояние (красная/зелёная) (результирующий - человек, машина, вход, задержка)
	TimeLasErr time.Time `json:"-"`               // Время, когда последний раз на весовой было нарушено соблюдение хотя-бы одного условия
}

type Cam struct {
	Id            string            `json:"-"`      // ИД-Камеры. Получаем по RESTу на основании serial, пользователю в конфиге он не нужен
	Serial        string            `json:"serial"` // Серийный номер камеры, по нему ёё и идентифицируем и заполняем её ID
	ConState      string            `json:"-"`      // Статус, получаем через WebPOint, например 'CONNECTED'
	Name          string            `json:"-"`      // Имя камеры, получаем актуальное через WebPOint
	Inputs        map[string]*Input `json:"-"`      // Состояние входов
	Car           bool              `json:"-"`      // В зоне обнаружена машина
	CarEventId    string            `json:"-"`      // Id события когда машина заехала в зону
	Person        bool              `json:"-"`      // В зоне обнаружен человек
	PersonEventId string            `json:"-"`      // Id события, когда человек зашел в зону
}

type Input struct {
	EntityId string `json:"-"` // Заполняем динамически, берём у камеры в links []{ {type:"DIGITAL_INPUT", source: "4xIx1DMwMLSwMDW2tDBKNNBLycwzMBASCDilIfJR0W3apqrIovO_tncAAA"},{},...}
	State    bool   `json:"-"` // Статус входа
}

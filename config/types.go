package config

import (
	"sync"
	"time"
)

type States string
type CamState string

const (
	ZoneNameAppendix string   = "ZoneNameIs"
	StateTrue        States   = "green"
	StateFalse       States   = "red"
	StateUnknown     States   = ""
	CamConnected     CamState = "CONNECTED"
	//CamFactoryDefault CamState = "FACTORY_DEFAULT" - Оказалось это плохой случай :)
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
	Nonce             string      `json:"-"`                   // Avigilon developer nonce
	DevNonce          string      `json:"dev_nonce,omitempty"` // DevNonce if we need it
	Key               string      `json:"-"`                   // Avigilon developer key
	DevKey            string      `json:"dev_key,omitempty"`   // DevKey if we need it
	Zones             []Zone      `json:"zones"`               // Список весовых зон
	Debug             bool        `json:"debug,omitempty"`     // Работа в режиме отладки
	mu                *sync.Mutex `json:"-"`
}

type Zone struct {
	Id             string    `json:"Id"`                // По нему будем работать с Зоной.?Zone=hexEncoded(ZoneNameAppendix+name) (Обязательно обновлять и сохранять в конфиге если при чтении конфига была пустая)
	Name           string    `json:"name"`              // Имя зоны -> в вебе будет использоваться для отображения (?Zone=hexEncoded(ZoneNameAppendix+name))
	Cameras        []Cam     `json:"cameras"`           // Камеры в пределах текущей зоны
	DelaySec       int       `json:"delay_sec"`         // Задержка после сработки входа, наличия машины и отсутствия человека
	TimeLasErr     time.Time `json:"-"`                 // Время, когда последний раз на весовой было нарушено соблюдение хотя-бы одного условия
	CarOnAnyCamera bool      `json:"car_on_any_camera"` // Если машина хоть на одной камере - всё ок
	IgnoreCarState bool      `json:"ignore_car_state"`  // Игнорировать состояние машины. Не сбрасывать таймер.
}

type Cam struct {
	Id             string            `json:"-"`                         // ИД-Камеры. Получаем по RESTу на основании serial, пользователю в конфиге он не нужен
	Serial         string            `json:"serial"`                    // Серийный номер камеры, по нему ёё и идентифицируем и заполняем её ID
	ConState       CamState          `json:"-"`                         // Статус, получаем через WebPoint, например 'CONNECTED'.
	Name           string            `json:"-"`                         // Имя камеры, получаем актуальное через WebPOint
	Inputs         map[string]*Input `json:"-"`                         // Состояние входов
	Car            States            `json:"-"`                         // В зоне обнаружена машина. "","green","red"
	CarEventId     string            `json:"-"`                         // Id события когда машина заехала в зону
	Person         States            `json:"-"`                         // В зоне обнаружен человек. "","green","red"
	PersonEventId  string            `json:"-"`                         // Id события, когда человек зашел в зону
	InputsDisabled bool              `json:"inputs_disabled,omitempty"` // Не использовать входы камеры
}

type Input struct {
	EntityId string `json:"-"` // Заполняем динамически, берём у камеры в links []{ {type:"DIGITAL_INPUT", source: "4xIx1DMwMLSwMDW2tDBKNNBLycwzMBASCDilIfJR0W3apqrIovO_tncAAA"},{},...}
	State    States `json:"-"` // Статус входа.  "","green","red"
}

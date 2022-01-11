package config

//certifiates for Avigilon Control Center Web Endpoint
//openssl req -new -x509 -days 365 -sha256 -newkey rsa:2048 -nodes -out avigilon.pem -keyout avigilon.key
//openssl req -new -x509 -days 365 -sha256 -newkey rsa:2048 -nodes -out public.crt -keyout private.key
//CN = Avigilon Control Center Web Endpoint
//OU = ACC
//O = Avigilon Corp.
//L = Vancouver
//S = British Columbia
//C = CA
type Config struct {
	WWWAddr           string `json:"www_addr"`            // Адрес на котором принимаем запросы
	WWWPort           string `json:"www_Port"`            // Порт на котором принимаем запросы
	WWWCertificate    string `json:"www_certificate"`     // Путь к файлу сертификатов https
	WWWCertificateKey string `json:"www_certificate_key"` // Путь к файлу ключа сертификата
	WPServer          string `json:"wp_server"`           // webPoint  адрес
	WPPort            string `json:"wp_port"`             // webPoint  порт
	WPUser            string `json:"wp_user"`             // webPoint  user
	WPPassword        string `json:"wp_password"`         // webPoint  password
	DevNonce          string `json:"dev_nonce,omitempty"` // Avigilon developer nonce
	DevKey            string `json:"dev_key,omitempty"`   // Avigilon developer key
	Zones             []zone `json:"zones"`               // Список весовых зон
	Debug             bool   `json:"debug,omitempty"`     // Работа в режиме отладки
}

type zone struct {
	Name      string `json:"name"`                // Имя зоны
	Cameras   []cam  `json:"cameras"`             // Камеры в пределах текущей зоны
	DelaySec  int    `json:"delay_sec"`           // Задержка после сработки входа, наличия машины и отсутствия человека
	Bookmarks bool   `json:"bookmarks,omitempty"` // Генерировать Закладки
	Alarms    bool   `json:"alarms,omitempty"`    // Генерировать тревоги
	State     bool   `json:"state,omitempty"`     // Текущее состояние (красная/зелёная) (результирующий - человек, машина, вход, задержка)
}

type cam struct {
	CamID    string  `json:"-"`                // ИД-Камеры. Получаем по RESTу на основании serial
	Serial   string  `json:"serial"`           // Серийный номер камеры, по нему и ёё и идентифицируем
	ConState string  `json:"connectionState"`  // 'CONNECTED'
	Name     string  `json:"name"`             // Имя камеры
	Inputs   []input `json:"inputs"`           // Состояние входов
	Car      bool    `json:"car,omitempty"`    // В зоне обнаружена машина
	Person   bool    `json:"person,omitempty"` // В зоне обнаружен человек
}

type input struct {
	Id       string `json:"id"`                 // НЕ уверен, надо проверить что приходит из веб-поинта.
	Name     string `json:"name"`               // Имя будет отображаться в вебе
	Inverion bool   `json:"inverion,omitempty"` // Если состояние входа необходимо инвертировать
	State    bool   `json:"state,omitempty"`    // Статус входа
}

func New() *Config {
	return &Config{}
}

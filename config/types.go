package config

type Config struct {
	WWWAddr           string `json:"www_addr"`            //Адрес на котором принимаем запросы
	WWWCertificate    string `json:"www_certificate"`     //путь к файлу сертификатов https
	WWWCertificateKey string `json:"www_certificate_key"` //путь к файлу ключа сетрификата
	WPServer          string `json:"wp_server"`           //webPoint  адрес
	WPPort            string `json:"wp_port"`             //webPoint  порт
	WPUser            string `json:"wp_user"`             //webPoint  user
	WPPassword        string `json:"wp_password"`         //webPoint  password
	DevNonce          string `json:"dev_nonce"`           //Avigilon developer nonce
	DevKey            string `json:"dev_key"`             //Avigilon developer key
	Zones             []zone `json:"zones"`               //список весовых зон
	Debug             bool   `json:"debug"`
}

type zone struct {
	Name      string `json:"name"`      //Имя зоны
	Cameras   []Cam  `json:"cameras"`   //Камеры
	DelaySec  int    `json:"delay_sec"` //Задержка после сработки входа, наличия машины и отсуцтвия человека
	Bookmarks bool   `json:"bookmarks"` //Генерировать Закладки
	Alarms    bool   `json:"alarms"`    //Генерировать тревоги
	State     bool   `json:"state"`     //текущее состояние (красная/зелёная)
}

type Cam struct {
	CamID  string  `json:"cam_id"` //Камера которую смотрим
	Name   string  `json:"name"`   //Имя камеры
	Inputs []input `json:"inputs"` //Состояние входов
	Car    bool    `json:"car"`    //В зоне обнаружена машина
	Person bool    `json:"person"` //В зоне обнаружен человек

}

type input struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State bool   `json:"state"`
}

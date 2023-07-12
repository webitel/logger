package model

type AppConfig struct {
	Rabbit   *RabbitConfig   `json:"rabbit,omitempty"`
	Database *DatabaseConfig `json:"database,omitempty"`
}

type RabbitConfig struct {
	Url string `json:"url"`
}

type DatabaseConfig struct {
	Url string `json:"url"`
}

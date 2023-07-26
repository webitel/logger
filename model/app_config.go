package model

type AppConfig struct {
	Rabbit   *RabbitConfig   `json:"rabbit,omitempty"`
	Database *DatabaseConfig `json:"database,omitempty"`
	Consul   *ConsulConfig   `json:"consul,omitempty"`
}

type RabbitConfig struct {
	Url string `json:"url"`
}

type DatabaseConfig struct {
	Url string `json:"url"`
}

type ConsulConfig struct {
	Id            string `json:"id"`
	Address       string `json:"address"`
	PublicAddress string `json:"publicAddress"`
}

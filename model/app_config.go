package model

type AppConfig struct {
	Rabbit   *RabbitConfig   `json:"rabbit,omitempty"`
	Database *DatabaseConfig `json:"database,omitempty"`
	Consul   *ConsulConfig   `json:"consul,omitempty"`
	Grpc     *GrpcConfig     `json:"grpc,omitempty"`
}

type RabbitConfig struct {
	Url string `json:"url"`
}

type DatabaseConfig struct {
	Url string `json:"url"`
}

type ConsulConfig struct {
	Id      string `json:"id"`
	Address string `json:"address"`
}

type GrpcConfig struct {
	Address string `json:"address"`
}

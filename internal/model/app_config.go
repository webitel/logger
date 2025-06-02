package model

type AppConfig struct {
	Rabbit   *RabbitConfig   `json:"rabbit,omitempty"`
	Database *DatabaseConfig `json:"database,omitempty"`
	Consul   *ConsulConfig   `json:"consul,omitempty"`
	Features *FeaturesConfig `json:"features,omitempty"`
	GrpcAddr string          `json:"publicAddress" flag:"grpc_addr|| Public grpc address with port" env:"GRPC_ADDR"`
}

type RabbitConfig struct {
	Url string `json:"url" flag:"amqp|| AMQP connection" env:"BROKER_ADDRESS"`
}

type DatabaseConfig struct {
	Url string `json:"url" flag:"data_source|| Data source" env:"DATASOURCE"`
}

type ConsulConfig struct {
	Id      string `json:"id" flag:"id|1| Service tag" env:"SERVICE_ID"`
	Address string `json:"address" flag:"consul|| Host to consul" env:"CONSUL"`
}

type FeaturesConfig struct {
	EnableLoginConsumption bool `json:"loginConsumption" flag:"login_consumption|true| Enables login consumption" env:"LOGIN_CONSUMPTION"`
	EnableCrudEvents       bool `json:"crudEvents" flag:"saga_events|false| Enable saga events" env:"CRUD_EVENTS"`
}

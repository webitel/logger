package model

type AppConfig struct {
	Rabbit   *RabbitConfig        `json:"rabbit,omitempty"`
	Database *DatabaseConfig      `json:"database,omitempty"`
	Consul   *ConsulConfig        `json:"consul,omitempty"`
	Log      *ObservabilityConfig `json:"log,omitempty" `
}

type RabbitConfig struct {
	Url string `json:"url" flag:"amqp|| AMQP connection"`
}

type DatabaseConfig struct {
	Url string `json:"url" flag:"data_source|| Data source"`
}

type ConsulConfig struct {
	Id            string `json:"id" flag:"id|1| Service tag" `
	Address       string `json:"address" flag:"consul|| Host to consul"`
	PublicAddress string `json:"publicAddress" flag:"grpc_addr|| Public grpc address with port"`
}

type ObservabilityConfig struct {
	SdkExport bool   `json:"sdk_export" flag:"otel_export|false| Export logs, metrics, traces to otlp or file. Configures by default env variables"`
	LogLevel  string `json:"log_lvl" flag:"log_lvl|debug| Max log level"`
}

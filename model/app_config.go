package model

const (
	ServiceVersion = "24.08"
	ServiceName    = "logger"
)

type AppConfig struct {
	Rabbit       RabbitConfig   `json:"rabbit,omitempty"`
	Database     DatabaseConfig `json:"database,omitempty"`
	Consul       ConsulConfig   `json:"consul,omitempty" `
	TracerConfig TracesConfig   `json:"traces,omitempty"`
	LoggerConfig LogsConfig     `json:"logs,omitempty"`
}

func (t *AppConfig) Normalize() AppError {
	err := t.Consul.Normalize()
	if err != nil {
		return err
	}
	err = t.TracerConfig.Normalize()
	if err != nil {
		return err
	}
	err = t.TracerConfig.Normalize()
	if err != nil {
		return err
	}
	return nil
}

type RabbitConfig struct {
	Url *string `json:"url" flag:"amqp|| AMQP connection"`
}

func (t *RabbitConfig) Normalize() AppError {
	if t.Url == nil {
		url := "amqp://admin:admin@localhost:5672?heartbeat=10"
		t.Url = &url
	}

	return nil
}

type DatabaseConfig struct {
	Url *string
}

func (t *DatabaseConfig) Normalize() AppError {
	if t.Url == nil {
		url := "postgres://postgres:postgres@localhost:5432/webitel"
		t.Url = &url
	}

	return nil
}

type ConsulConfig struct {
	Id            *string
	Address       *string
	PublicAddress *string
}

func (t *ConsulConfig) Normalize() AppError {
	if t.Id == nil {
		id := "logger"
		t.Id = &id
	}
	if t.Address == nil {
		a := "127.0.0.1:8500"
		t.Address = &a
	}
	if t.PublicAddress == nil {
		a := "127.0.0.1:10001"
		t.PublicAddress = &a
	}
	return nil
}

type TracesConfig struct {
	Provider *string /*`json:"provider" flag:"tracer_provider|t| Collector's type (otlp|stdout|jaeger)"`*/
	Address  *string /*`json:"address" flag:"tracer_address|t| Connection to the tracer collector endpoint if needed (format x.x.x.x:xxxx)"`*/
}

func (t *TracesConfig) Normalize() AppError {
	defaultProvider := "stdout"
	if t.Provider == nil || *t.Provider == "" {
		t.Provider = &defaultProvider
	}
	if t.Address == nil {
		a := ""
		t.Address = &a
	}
	return nil
}

//	type MeterConfig struct {
//		Provider string `json:"provider" flag:"meter_provider|stdout| Collector's type (otlp|stdout|jaeger)"`
//		Address  string `json:"address" flag:"meter_address|127.0.0.1:4317| Connection to the collector endpoint if needed(format x.x.x.x:xxxx)"`
//	}
type LogsConfig struct {
	Provider *string
	Address  *string
	LogLevel *string
}

func (t *LogsConfig) Normalize() AppError {
	defaultProvider := "wlog"
	if t.Provider == nil || *t.Provider == "" {
		t.Provider = &defaultProvider
	}
	if t.Address == nil {
		a := ""
		t.Address = &a
	}
	if t.LogLevel == nil {
		a := "info"
		t.LogLevel = &a
	}
	return nil
}

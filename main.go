package main

import (
	"flag"
	"fmt"
	"github.com/webitel/logger/app"
	"github.com/webitel/logger/model"
	"github.com/webitel/webitel-go-kit/logs"
	"github.com/webitel/wlog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	configPath *string
)

func main() {
	config, appErr := loadConfig()
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}
	appErr = config.Normalize()
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}

	// init of logger
	externalLogger, appErr := initLogger(&config.LoggerConfig)
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}
	// Create an application layer
	app, appErr := app.New(config, externalLogger)
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}
	initSignals(app)
	appErr = app.Start()
	wlog.Critical(appErr.Error())
	return

}

func initSignals(app *app.App) {
	wlog.Info("initializing stop signals")
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)

	go func() {
		for {
			s := <-sigchnl
			handleSignals(s, app)
		}
	}()

}

func handleSignals(signal os.Signal, app *app.App) {
	if signal == syscall.SIGTERM || signal == syscall.SIGINT || signal == syscall.SIGKILL {
		app.Stop()
		wlog.Info(fmt.Sprintf("got kill signal, service gracefully stopped!"))

		os.Exit(0)
	}
}

func loadConfig() (*model.AppConfig, model.AppError) {
	var appConfig model.AppConfig

	// Tracer
	appConfig.TracerConfig.Provider = flag.String("tracing_provider", "", "Collector's type (otlp|stdout|jaeger)")
	appConfig.TracerConfig.Address = flag.String("tracing_address", "", "Connection to the tracer collector endpoint if needed (format x.x.x.x:xxxx)")
	// Logger
	appConfig.LoggerConfig.Provider = flag.String("logging_provider", "", "Collector's type (otlp|stdout|wlog)")
	appConfig.LoggerConfig.Address = flag.String("logging_address", "", "Connection to the logger collector endpoint if needed (format x.x.x.x:xxxx)")
	appConfig.LoggerConfig.LogLevel = flag.String("logging_level", "info", "Connection to the logger collector endpoint if needed (format x.x.x.x:xxxx)")

	// Consul
	appConfig.Consul.Id = flag.String("id", "1", "Service tag")
	appConfig.Consul.Address = flag.String("consul", "", "Host to consul")
	appConfig.Consul.PublicAddress = flag.String("grpc_addr", "", "Public grpc address with port")
	// Database
	appConfig.Database.Url = flag.String("data_source", "", "Data source")
	// Rabbit
	appConfig.Rabbit.Url = flag.String("amqp", "", "AMQP connection")

	flag.Parse()

	return &appConfig, nil
}

func initLogger(conf *model.LogsConfig) (*logs.DynamicLogger, model.AppError) {
	if conf == nil {
		p := "wlog"
		conf = &model.LogsConfig{Provider: &p}
	}
	appErr := conf.Normalize()
	if appErr != nil {
		return nil, appErr
	}

	defaultOut := wlog.NewLogger(&wlog.LoggerConfiguration{EnableConsole: true, ConsoleLevel: wlog.LevelDebug})

	l, err := logs.New(
		defaultOut,
		logs.WithTimeout(time.Second*5),
		logs.WithAddress(*conf.Address),
		logs.WithServiceName(model.ServiceName),
		logs.WithServiceVersion(model.ServiceVersion),
		logs.WithExporter(*conf.Provider),
		logs.WithLogLevel(*conf.LogLevel),
	)

	if err != nil {
		return nil, model.NewInternalError("app.app.init_logger.default.create.error", err.Error())
	}
	return l, nil
}

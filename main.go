package main

import (
	"context"
	"fmt"
	"github.com/BoRuDar/configuration/v4"
	"github.com/webitel/logger/app"
	"github.com/webitel/logger/model"
	otelsdk "github.com/webitel/webitel-go-kit/otel/sdk"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	// -------------------- plugin(s) -------------------- //
	_ "github.com/webitel/webitel-go-kit/otel/sdk/log/otlp"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/log/stdout"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/metric/otlp"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/metric/stdout"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/trace/otlp"
	_ "github.com/webitel/webitel-go-kit/otel/sdk/trace/stdout"
)

var (
	name      = "logger"
	version   = "24.08"
	namespace = "webitel"
	serviceId string
)

func main() {

	config, appErr := loadConfig()
	if appErr != nil {
		slog.Error(appErr.Error())
		return
	}
	serviceId = config.Consul.Id

	if config.Log.SdkExport {
		sd := SetupOtel(config.Log.LogLevel)
		defer sd(context.Background())
	} else {
		var lvl slog.Level
		switch config.Log.LogLevel {
		case "info":
			lvl = slog.LevelInfo
		case "warn":
			lvl = slog.LevelWarn
		case "error":
			lvl = slog.LevelError
		case "debug":
			lvl = slog.LevelDebug
		default:
			slog.Info("unable to determine log level, setting debug")
			lvl = slog.LevelDebug
		}
		slog.SetLogLoggerLevel(lvl)
	}

	// * Create an application layer
	app, appErr := app.New(config)
	if appErr != nil {
		slog.Error(appErr.Error())
		return
	}
	initSignals(app)
	appErr = app.Start()
	slog.Error(appErr.Error())
	return

}

func SetupOtel(severity string) (shutdown otelsdk.ShutdownFunc) {

	service := resource.NewSchemaless(
		semconv.ServiceName(name),
		semconv.ServiceVersion(version),
		semconv.ServiceInstanceID(serviceId),
		semconv.ServiceNamespace(namespace),
	)
	ctx := context.Background()

	var lvl log.Severity

	switch severity {
	case "info":
		lvl = log.SeverityInfo
	case "warn":
		lvl = log.SeverityWarn
	case "error":
		lvl = log.SeverityError
	case "debug":
		lvl = log.SeverityDebug
	default:
		slog.Info("unable to determine log level, setting debug")
		lvl = log.SeverityDebug
	}
	shutdown, err := otelsdk.Setup(ctx, otelsdk.WithResource(service), otelsdk.WithLogLevel(lvl))
	if err != nil {
		slog.Error(err.Error())
		return
	}
	stdlog := otelslog.NewLogger(name)
	slog.SetDefault(stdlog)
	return shutdown
}

func initSignals(app *app.App) {
	slog.Info("initializing stop signals")
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
		slog.Info(fmt.Sprintf("got kill signal, service gracefully stopped!"))

		os.Exit(0)
	}
}

func loadConfig() (*model.AppConfig, model.AppError) {
	var appConfig model.AppConfig

	configurator := configuration.New(
		&appConfig,
		// order of execution will be preserved:
		configuration.NewFlagProvider(),
		configuration.NewEnvProvider(),
		configuration.NewDefaultProvider(),
	).SetOptions(
		configuration.OnFailFnOpt(func(err error) {
			//fmt.Printf(err.Error())
		}))

	if err := configurator.InitValues(); err != nil {
		return nil, model.NewInternalError("main.main.unmarshal_config.bad_arguments.parse_fail", err.Error())
	}
	return &appConfig, nil
}

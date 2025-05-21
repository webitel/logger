package logger

import (
	"context"
	"fmt"
	"github.com/BoRuDar/configuration/v4"
	"github.com/webitel/logger/internal/app"
	"github.com/webitel/logger/internal/model"
	otelsdk "github.com/webitel/webitel-go-kit/otel/sdk"
	"go.opentelemetry.io/contrib/bridges/otelslog"
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

const (
	name      = "logger"
	version   = "25.04"
	namespace = "webitel"
)

func Run() error {

	config, appErr := loadConfig()
	if appErr != nil {
		return appErr
	}
	sd := SetupOtel(config.Consul.Id)
	defer sd(context.Background())
	// * Create an application layer
	app, err := app.New(config)
	if err != nil {
		return err
	}
	initSignals(app)
	err = app.Start()
	return err

}

func SetupOtel(serviceId string) (shutdown otelsdk.ShutdownFunc) {

	service := resource.NewSchemaless(
		semconv.ServiceName(name),
		semconv.ServiceVersion(version),
		semconv.ServiceInstanceID(serviceId),
		semconv.ServiceNamespace(namespace),
	)
	shutdown, err := otelsdk.Configure(context.Background(), otelsdk.WithResource(service),
		otelsdk.WithLogBridge(
			func() {
				// Create new slog logger with otel custom log level and handler
				logger := slog.New(
					otelslog.NewHandler("slog"),
				)
				slog.SetDefault(logger)
			},
		),
	)
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

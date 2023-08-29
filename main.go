package main

import (
	"flag"
	"fmt"
	"github.com/BoRuDar/configuration/v4"
	"github.com/webitel/logger/api"
	"github.com/webitel/logger/app"
	"github.com/webitel/logger/model"
	"github.com/webitel/logger/rabbit"
	"github.com/webitel/logger/storage"
	"github.com/webitel/logger/storage/postgres"
	"github.com/webitel/wlog"
	"os"
	"os/signal"
	"syscall"

	errors "github.com/webitel/engine/model"
)

var (
	configPath *string
)

func main() {
	//flagDefine()
	log := wlog.NewLogger(&wlog.LoggerConfiguration{
		EnableConsole: true,
		ConsoleLevel:  wlog.LevelDebug,
	})

	wlog.RedirectStdLog(log)
	wlog.InitGlobalLogger(log)

	config, appErr := UnmarshalConfig()
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}
	store, appErr := BuildDatabase(config.Database)
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}
	defer store.Close()
	initSignals(store)
	// * Create an application layer
	app, appErr := app.New(store, config)
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}

	errChan := make(chan errors.AppError)
	// * Build and run rabbit listener
	go rabbit.BuildAndServe(app, config.Rabbit, errChan)
	// * Build and run grpc server
	go api.ServeRequests(app, config.Consul, errChan)

	appErr = <-errChan
	wlog.Critical(appErr.Error())
	return

}

func flagDefine() {
	configPath = flag.String("config", "./config/config.json", "Path to the config file")
	flag.Parse()
}

func BuildDatabase(config *model.DatabaseConfig) (storage.Storage, errors.AppError) {
	store, err := postgres.New(config)
	if err != nil {
		return nil, err
	}
	err = store.Open()
	if err != nil {
		return nil, err
	}
	//err = store.SchemaInit()
	//if err != nil {
	//	return nil, err
	//}
	return store, nil
}

func initSignals(store storage.Storage) {
	wlog.Info("initializing stop signals")
	sigchnl := make(chan os.Signal, 1)
	signal.Notify(sigchnl)

	go func() {
		for {
			s := <-sigchnl
			handleSignals(s, store)
		}
	}()

}

func handleSignals(signal os.Signal, store storage.Storage) {
	if signal == syscall.SIGTERM || signal == syscall.SIGINT || signal == syscall.SIGKILL {
		wlog.Info(fmt.Sprintf("got kill signal. service gracefully stopped!"))
		store.Close()
		os.Exit(0)
	}
}

func UnmarshalConfig() (*model.AppConfig, errors.AppError) {
	var appConfig model.AppConfig
	//if configPath == nil {
	//	return nil, errors.NewBadRequestError("main.main.unmarshal_config.bad_arguments.config_path_is_nil", "config path is nil")
	//}

	configurator := configuration.New(
		&appConfig,
		// order of execution will be preserved:
		configuration.NewFlagProvider(),
		configuration.NewEnvProvider(),
		configuration.NewDefaultProvider(),
	)

	if err := configurator.InitValues(); err != nil {
		return nil, errors.NewInternalError("main.main.unmarshal_config.bad_arguments.parse_fail", err.Error())
	}

	//file, err := ioutil.ReadFile(*configPath)
	//if err != nil {
	//	return nil, errors.NewBadRequestError("main.main.unmarshal_config.bad_arguments.wrong_path", err.Error())
	//}
	//err = json.Unmarshal(file, &appConfig)
	//if err != nil {
	//	return nil, errors.NewBadRequestError("main.main.unmarshal_config.json_unmarshal.fail", err.Error())
	//}
	return &appConfig, nil
}

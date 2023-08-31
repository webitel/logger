package main

import (
	"fmt"
	"github.com/BoRuDar/configuration/v4"
	"github.com/webitel/logger/app"
	"github.com/webitel/logger/model"
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

	config, appErr := loadConfig()
	if appErr != nil {
		wlog.Critical(appErr.Error())
		return
	}
	//store, appErr := BuildDatabase(config.Database)
	//if appErr != nil {
	//	wlog.Critical(appErr.Error())
	//	return
	//}
	//defer store.Close()
	//initSignals(store)
	// * Create an application layer
	app, appErr := app.New(config)
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
		wlog.Info(fmt.Sprintf("got kill signal, service gracefully stopped!"))
		app.Stop()
		os.Exit(0)
	}
}

func loadConfig() (*model.AppConfig, errors.AppError) {
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

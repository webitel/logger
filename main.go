package main

import (
	"encoding/json"
	"flag"
	"github.com/webitel/wlog"
	"io/ioutil"
	"log"
	"logger/api"
	"logger/app"
	"logger/model"
	"logger/rabbit"
	"logger/storage"
	"logger/storage/postgres"
	"os"
	"os/signal"
	"syscall"

	errors "github.com/webitel/engine/model"
)

var (
	configPath *string
)

func main() {
	flagDefine()

	config, appErr := UnmarshalConfig()
	if appErr != nil {
		log.Fatal(appErr.Error())
	}
	store, appErr := BuildDatabase(config.Database)
	if appErr != nil {
		log.Fatal(appErr.Error())
	}
	defer store.Close()
	initSignals(store)
	// * Create an application layer
	app, appErr := app.New(store, config)
	if appErr != nil {
		log.Fatal(appErr.Error())
	}

	errChan := make(chan errors.AppError)
	// * Build and run rabbit listener
	go rabbit.BuildAndServe(app, config.Rabbit, errChan)
	// * Build and run grpc server
	go api.ServeRequests(app, config.Consul, errChan)
	log.Fatal(<-errChan)

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
		wlog.Info("got kill signal. ")
		wlog.Info("program will terminate now.")
		store.Close()
		os.Exit(0)
	}
}

func UnmarshalConfig() (*model.AppConfig, errors.AppError) {
	var appConfig model.AppConfig
	if configPath == nil {
		return nil, errors.NewBadRequestError("main.main.unmarshal_config.bad_arguments.config_path_is_nil", "config path is nil")
	}

	file, err := ioutil.ReadFile(*configPath)
	if err != nil {
		return nil, errors.NewBadRequestError("main.main.unmarshal_config.bad_arguments.wrong_path", err.Error())
	}
	err = json.Unmarshal(file, &appConfig)
	if err != nil {
		return nil, errors.NewBadRequestError("main.main.unmarshal_config.json_unmarshal.fail", err.Error())
	}
	return &appConfig, nil
}

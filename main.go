package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"webitel_logger/api"
	"webitel_logger/app"
	"webitel_logger/model"
	"webitel_logger/proto"
	"webitel_logger/storage"
	"webitel_logger/storage/postgres"

	errors "github.com/webitel/engine/model"
	"google.golang.org/grpc"
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

	server, appErr := BuildGrpc(store)
	if appErr != nil {
		log.Fatal(appErr.Error())
	}

	//  * Open tcp connection
	listener, err := net.Listen("tcp", "localhost:1137")
	if err != nil {
		log.Fatal(err)
	}
	err = server.Serve(listener)
	if err != nil {
		log.Fatal()
	}

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
	err = store.SchemaInit()
	if err != nil {
		return nil, err
	}
	return store, nil
}

func BuildGrpc(store storage.Storage) (*grpc.Server, errors.AppError) {

	// * Create an application layer
	app, err := app.New(store)
	if err != nil {
		log.Fatal(err.Error())
	}
	grpcServer := grpc.NewServer()
	// * Creating services
	l, appErr := api.NewLoggerService(app)
	if appErr != nil {
		return nil, appErr
	}
	c, appErr := api.NewConfigService(app)
	if appErr != nil {
		return nil, appErr
	}

	// * register logger service
	proto.RegisterLoggerServiceServer(grpcServer, l)
	// * register config service
	proto.RegisterConfigServiceServer(grpcServer, c)

	return grpcServer, nil
}

func BuildRabbit(config storage.Storage) errors.AppError {

	return nil
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

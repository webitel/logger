package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"webitel_logger/api"
	"webitel_logger/app"
	"webitel_logger/model"
	"webitel_logger/proto"
	"webitel_logger/rabbit"
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
	initSignals(store)
	// * Create an application layer
	app, appErr := app.New(store)
	if appErr != nil {
		log.Fatal(appErr.Error())
	}

	errChan := make(chan errors.AppError)
	// * Build and run rabbit listener
	go BuildRabbit(app, config.Rabbit, errChan)
	// * Build and run grpc server
	go ServeRequests(store, app, errChan)
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
	err = store.SchemaInit()
	if err != nil {
		return nil, err
	}
	return store, nil
}

func BuildGrpc(store storage.Storage, app *app.App) (*grpc.Server, errors.AppError) {

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

func BuildRabbit(app *app.App, config *model.RabbitConfig, errChan chan errors.AppError) {
	handler, err := rabbit.NewHandler(app)
	if err != nil {
		errChan <- err
		return
	}
	listener, err := rabbit.NewListener(config, errChan)
	if err != nil {
		errChan <- err
		return
	}

	listener.Listen(handler.Handle)
}

func ServeRequests(store storage.Storage, app *app.App, errChan chan errors.AppError) {
	// * Build grpc server
	server, appErr := BuildGrpc(store, app)
	if appErr != nil {
		errChan <- appErr
		return
	}
	//  * Open tcp connection
	listener, err := net.Listen("tcp", "localhost:1137")
	if err != nil {
		errChan <- errors.NewInternalError("main.main.serve_requests.listen.error", err.Error())
		return
	}
	err = server.Serve(listener)
	if err != nil {
		errChan <- errors.NewInternalError("main.main.serve_requests.serve.error", err.Error())
		return
	}
}

func initSignals(store storage.Storage) {
	log.Println("initializing stop signals")
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
		log.Println("got kill signal. ")
		log.Println("program will terminate now.")
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

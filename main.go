package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"webitel_logger/api"
	"webitel_logger/app"
	"webitel_logger/model"
	"webitel_logger/rabbit"
	"webitel_logger/storage"
	"webitel_logger/storage/postgres"

	"github.com/webitel/engine/discovery"
	errors "github.com/webitel/engine/model"
)

var (
	configPath                  *string
	address                     string
	APP_SERVICE_TTL             = time.Second * 30
	APP_DEREGESTER_CRITICAL_TTL = time.Minute * 2
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
	appErr = ConnectConsul(config.Consul)
	if appErr != nil {
		log.Fatal(appErr.Error())
	}

	errChan := make(chan errors.AppError)
	// * Build and run rabbit listener
	go rabbit.BuildAndServe(app, config.Rabbit, errChan)
	// * Build and run grpc server
	go ServeRequests(store, app, errChan)
	log.Fatal(<-errChan)

}

func flagDefine() {
	configPath = flag.String("config", "./config/config.json", "Path to the config file")
	flag.Parse()
}

func ConnectConsul(config *model.ConsulConfig) errors.AppError {
	consul, err := discovery.NewConsul(config.Id, config.Address, func() (bool, error) {
		return true, nil
	})
	err = consul.RegisterService("logger", address, 0, APP_SERVICE_TTL, APP_DEREGESTER_CRITICAL_TTL)
	if err != nil {
		return errors.NewInternalError("main.main.build_consul.register_in_consul.error", err.Error())
	}
	return nil
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

// func BuildGrpc(store storage.Storage, app *app.App) (*grpc.Server, errors.AppError) {

// 	grpcServer := grpc.NewServer()
// 	// * Creating services
// 	l, appErr := api.NewLoggerService(app)
// 	if appErr != nil {
// 		return nil, appErr
// 	}
// 	c, appErr := api.NewConfigService(app)
// 	if appErr != nil {
// 		return nil, appErr
// 	}

// 	// * register logger service
// 	proto.RegisterLoggerServiceServer(grpcServer, l)
// 	// * register config service
// 	proto.RegisterConfigServiceServer(grpcServer, c)

// 	return grpcServer, nil
// }

// func ServeRabbit(app *app.App, config *model.RabbitConfig, errChan chan errors.AppError) {
// 	handler, err := rabbit.NewHandler(app)
// 	if err != nil {
// 		errChan <- err
// 		return
// 	}
// 	listener, err := rabbit.NewListener(config, errChan)
// 	if err != nil {
// 		errChan <- err
// 		return
// 	}

// 	listener.Listen(handler.Handle)
// }

func ServeRequests(store storage.Storage, app *app.App, errChan chan errors.AppError) {
	// * Build grpc server
	server, appErr := api.BuildGrpc(store, app)
	if appErr != nil {
		errChan <- appErr
		return
	}
	//  * Open tcp connection
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:0", address))
	if err != nil {
		errChan <- errors.NewInternalError("main.main.serve_requests.listen.error", err.Error())
		return
	}

	fmt.Println(listener.Addr().String())

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
	address = appConfig.Grpc.Address
	return &appConfig, nil
}

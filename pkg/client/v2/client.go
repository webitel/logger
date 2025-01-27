package client

import (
	loggerapi "buf.build/gen/go/webitel/logger/grpc/go/_gogrpc"
	loggerproto "buf.build/gen/go/webitel/logger/protocolbuffers/go"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cache "github.com/hashicorp/golang-lru/v2/expirable"
	_ "github.com/mbobakov/grpc-consul-resolver"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Action string

func (a Action) String() string {
	return string(a)
}

const (
	CreateAction Action = "create"
	UpdateAction Action = "update"
	DeleteAction Action = "delete"
	ReadAction   Action = "read"
)

const (
	ExchangeName            = "logger"
	ConsulLoggerServiceName = "logger"
	DefaultCacheTtl         = 120 * time.Second
)

var (
	LoggerClientGrpcErr = errors.New("can't check the object state")
)

type LoggerClient struct {
	grpcConnection  *grpc.ClientConn
	grpcClient      loggerapi.ConfigServiceClient
	memoryCache     *cache.LRU[string, bool]
	cacheTimeToLive time.Duration

	rabbitConnection *amqp.Connection
	channel          *amqp.Channel
}
type ObjectedLogger struct {
	objClass string
	parent   *LoggerClient
}

type LoggerClientOpts func(*LoggerClient) error

func WithCustomCacheTime(timeToLive time.Duration) LoggerClientOpts {
	return func(client *LoggerClient) error {
		client.cacheTimeToLive = timeToLive
		return nil
	}
}

func WithGrpcConnection(conn *grpc.ClientConn) LoggerClientOpts {
	return func(client *LoggerClient) error {
		if conn == nil {
			return errors.New("grpc connections required")
		}
		if conn.GetState() != connectivity.Ready {
			return errors.New("grpc connection should be opened")
		}
		client.grpcConnection = conn
		return nil
	}
}

func WithAmqpConnection(conn *amqp.Connection) LoggerClientOpts {
	return func(client *LoggerClient) error {
		if conn == nil {
			return errors.New("rabbit  connections required")
		}
		if conn.IsClosed() {
			return errors.New("rabbit connection should be opened")
		}
		client.rabbitConnection = conn
		return nil
	}
}

func WithGrpcConsulAddress(consulIpAddress string) LoggerClientOpts {
	return func(client *LoggerClient) error {
		conn, err := grpc.NewClient(fmt.Sprintf("consul://%s/%s?wait=14s", consulIpAddress, ConsulLoggerServiceName),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return err
		}
		client.grpcConnection = conn
		return nil
	}

}

func WithAmqpConnectionString(rabbitConnectionString string) LoggerClientOpts {
	return func(client *LoggerClient) error {
		conn, err := amqp.Dial(rabbitConnectionString)
		if err != nil {
			return err
		}
		client.rabbitConnection = conn
		return nil
	}

}

func NewLoggerClient(opts ...LoggerClientOpts) (*LoggerClient, error) {
	var err error
	logger := &LoggerClient{}
	for _, opt := range opts {
		err = opt(logger)
		if err != nil {
			return nil, err
		}
	}
	// validation
	if logger.cacheTimeToLive <= 0 {
		logger.cacheTimeToLive = DefaultCacheTtl
	}
	if logger.rabbitConnection == nil || logger.grpcConnection == nil {
		return nil, errors.New("rabbit and grpc connections required")
	}
	if logger.rabbitConnection.IsClosed() {
		return nil, errors.New("rabbit connection should be opened")
	}
	if logger.grpcConnection.GetState() == connectivity.Shutdown {
		return nil, errors.New("grpc connection should be opened")
	}
	// initialization
	logger.memoryCache = cache.NewLRU[string, bool](0, nil, logger.cacheTimeToLive)
	logger.channel, err = logger.rabbitConnection.Channel()
	if err != nil {
		return nil, err
	}
	err = logger.channel.ExchangeDeclare(
		ExchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, err
	}

	logger.grpcClient = loggerapi.NewConfigServiceClient(logger.grpcConnection)
	return logger, nil
}

func (l *LoggerClient) GetObjectedLogger(objClass string) *ObjectedLogger {
	return &ObjectedLogger{
		objClass: objClass,
		parent:   l,
	}
}

func (l *LoggerClient) sendContext(ctx context.Context, domainId int64, objclass string, message *Message) error {
	enabled, err := l.checkObjectConfig(ctx, domainId, objclass)
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = l.channel.PublishWithContext(ctx, ExchangeName, formatKey(domainId, objclass), false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return err
	}
	return nil
}

func (l *LoggerClient) checkObjectConfig(ctx context.Context, domainId int64, objclass string) (bool, error) {
	cacheKey := fmt.Sprintf("%d.%s", domainId, objclass)
	enabled, found := l.memoryCache.Get(cacheKey)
	if !found {
		resp, err := l.grpcClient.CheckConfigStatus(ctx, &loggerproto.CheckConfigStatusRequest{
			ObjectName: objclass,
			DomainId:   domainId,
		})
		if err != nil {
			return false, errors.Join(LoggerClientGrpcErr, err)
		}
		enabled = resp.GetIsEnabled()
		l.memoryCache.Add(cacheKey, enabled)
	}
	return enabled, nil
}

func (l *ObjectedLogger) SendContext(ctx context.Context, domainId int64, message *Message) error {
	return l.parent.sendContext(ctx, domainId, l.objClass, message)
}

func (l *ObjectedLogger) GetObjClass() string {
	return l.objClass
}

// region UTILITY

func formatKey(domainId int64, objClass string) string {
	return fmt.Sprintf("logger.%d.%s", domainId, objClass)
}

// endregion

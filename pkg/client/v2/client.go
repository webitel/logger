package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	cache "github.com/hashicorp/golang-lru/v2/expirable"
	_ "github.com/mbobakov/grpc-consul-resolver"
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
	publisher Publisher
	consumer  Consumer

	memoryCache     *cache.LRU[string, bool]
	cacheTimeToLive time.Duration
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

type Publisher interface {
	Publish(ctx context.Context, routingKey string, body []byte, headers map[string]any) error
}

type Consumer interface {
	Consume(ctx context.Context, routingKey string, handler func()) error
}

func WithPublisher(publisher Publisher) LoggerClientOpts {
	return func(client *LoggerClient) error {
		if publisher == nil {
			return errors.New("publisher is nil")
		}
		client.publisher = publisher
		return nil
	}
}

func WithConsumer(consumer Consumer) LoggerClientOpts {
	return func(client *LoggerClient) error {
		if consumer == nil {
			return errors.New("consumer is nil")
		}
		client.consumer = consumer
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
	if logger.publisher == nil {
		return nil, errors.New("publisher is nil")
	}
	if logger.consumer == nil {
		return nil, errors.New("consumer is nil")
	}
	// validation
	if logger.cacheTimeToLive <= 0 {
		logger.cacheTimeToLive = DefaultCacheTtl
	}
	// initialization
	logger.memoryCache = cache.NewLRU[string, bool](0, nil, logger.cacheTimeToLive)
	return logger, nil
}

func (l *LoggerClient) GetObjectedLogger(objClass string) *ObjectedLogger {
	return &ObjectedLogger{
		objClass: objClass,
		parent:   l,
	}
}

func (l *LoggerClient) sendContext(ctx context.Context, domainId int64, objclass string, message *Message) (operationId string, err error) {
	message.OperationId = uuid.NewString()
	enabled, err := l.checkObjectConfig(ctx, domainId, objclass)
	if err != nil {
		return operationId, err
	}
	if !enabled {
		return operationId, nil
	}
	body, err := json.Marshal(message)
	if err != nil {
		return operationId, err
	}
	err = l.publisher.Publish(ctx, formatKey(domainId, objclass), body, nil)
	if err != nil {
		return operationId, err
	}
	return operationId, nil
}

func (l *LoggerClient) checkObjectConfig(ctx context.Context, domainId int64, objclass string) (bool, error) {
	cacheKey := fmt.Sprintf("%d.%s", domainId, objclass)
	enabled, found := l.memoryCache.Get(cacheKey)
	if !found {
		l.memoryCache.Add(cacheKey, enabled)
	}
	return enabled, nil
}

func (l *ObjectedLogger) SendContext(ctx context.Context, domainId int64, message *Message) (operationId string, err error) {
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

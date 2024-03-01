package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	cache "github.com/hashicorp/golang-lru/v2/expirable"
	proto "github.com/webitel/logger/pkg/client/api"
)

const (
	DefaultCacheTimeout = 120 * time.Second
)

type GrpcClient interface {
	Start() error
	Stop()
	IsOpened() bool
	Config() ConfigApi
}

type grpcClient struct {
	consulAddress string
	connection    *grpc.ClientConn
	configClient  proto.ConfigServiceClient
	isOpened      bool
	memoryCache   *cache.LRU[string, bool]
	config        ConfigApi
}

func (c *grpcClient) IsOpened() bool {
	return c.isOpened
}

func (c *grpcClient) Start() error {
	conn, err := grpc.Dial(fmt.Sprintf("consul://%s/logger?wait=14s", c.consulAddress),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	client := proto.NewConfigServiceClient(c.connection)
	c.configClient = client
	c.isOpened = true
	c.connection = conn
	return nil
}

func (c *grpcClient) Stop() {
	c.connection.Close()
	c.configClient = nil
}

func (c *grpcClient) Config() ConfigApi {
	return c.config
}

func NewGrpcClient(consulAddr string) GrpcClient {
	client := &grpcClient{
		consulAddress: consulAddr,
	}
	client.config = NewConfigApi(client)
	client.memoryCache = cache.NewLRU[string, bool](200, nil, DefaultCacheTimeout)
	return client
}

type ConfigApi interface {
	CheckIsActive(ctx context.Context, domainId int64, objectName string) (bool, error)
}

func NewConfigApi(cli *grpcClient) ConfigApi {
	return &configApi{client: cli}
}

func FormatKey(domainId int64, objectName string) string {
	return fmt.Sprintf("logger.config.%d.%s", domainId, objectName)
}

type configApi struct {
	client *grpcClient
}

func (c *configApi) CheckIsActive(ctx context.Context, domainId int64, objectName string) (bool, error) {
	cacheKey := FormatKey(domainId, objectName)
	enabled, ok := c.client.memoryCache.Get(cacheKey)
	if !ok {
		in := &proto.CheckConfigStatusRequest{
			ObjectName: objectName,
			DomainId:   domainId,
		}
		res, err := c.client.configClient.CheckConfigStatus(ctx, in)
		if err != nil {
			return false, err
		}
		c.client.memoryCache.Add(cacheKey, res.GetIsEnabled())
		return res.GetIsEnabled(), nil
	}
	return enabled, nil
}

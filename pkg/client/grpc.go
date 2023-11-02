package client

import (
	"context"
	"fmt"

	"github.com/webitel/logger/pkg/cache"
	"github.com/webitel/logger/pkg/discovery"

	proto "github.com/webitel/protos/logger"
	"github.com/webitel/wlog"
)

const (
	WatcherInterval = 5 * 1000
)

type GrpcClient interface {
	Start() error
	Stop()
	IsOpened() bool
	Config() ConfigApi
}

type grpcClient struct {
	//stop             chan struct{}
	//stopped          chan struct{}
	serviceDiscovery discovery.ServiceDiscovery
	poolConnections  discovery.Pool
	//watcher          *discovery.Watcher
	//startOnce        sync.Once
	isOpened    bool
	memoryCache cache.CacheStore

	config ConfigApi
}

func (c *grpcClient) IsOpened() bool {
	return c.isOpened
}

func (c *grpcClient) Start() error {
	if services, err := c.serviceDiscovery.GetByName("logger"); err != nil {
		return err
	} else {
		for _, v := range services {
			c.registerConnection(v)
		}
	}
	c.isOpened = true
	return nil
}

func (c *grpcClient) Stop() {

	if c.poolConnections != nil {
		c.poolConnections.CloseAllConnections()
	}
	c.isOpened = false
	//close(c.stop)
	//<-c.stopped
}

func (c *grpcClient) registerConnection(v *discovery.ServiceConnection) {
	addr := fmt.Sprintf("%s:%d", v.Host, v.Port)
	client, err := NewConnection(v.Id, addr)
	if err != nil {
		wlog.Error(fmt.Sprintf("connection %s [%s] error: %s", v.Id, addr, err.Error()))
		return
	}
	c.poolConnections.Append(client)
	wlog.Debug(fmt.Sprintf("register connection %s [%s]", client.Name(), addr))
}

func (c *grpcClient) wakeUp() {
	list, err := c.serviceDiscovery.GetByName("logger")
	if err != nil {
		wlog.Error(err.Error())
		return
	}

	for _, v := range list {
		if _, err := c.poolConnections.GetById(v.Id); err == discovery.ErrNotFoundConnection {
			c.registerConnection(v)
		}
	}
	c.poolConnections.RecheckConnections(list.Ids())
}

func (c *grpcClient) getRandomClient() (*Connection, error) {
	cli, err := c.poolConnections.Get(discovery.StrategyRoundRobin)
	if err != nil {
		return nil, err
	}

	return cli.(*Connection), nil
}

// func (c *grpcClient) getClient(appId string) (*Connection, error) {
// 	cli, err := c.poolConnections.GetById(appId)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return cli.(*Connection), nil
// }

func (c *grpcClient) Config() ConfigApi {
	return c.config
}

func NewGrpcClient(serviceDiscovery discovery.ServiceDiscovery) GrpcClient {
	client := &grpcClient{
		//stop:             make(chan struct{}),
		//stopped:          make(chan struct{}),
		poolConnections:  discovery.NewPoolConnections(),
		serviceDiscovery: serviceDiscovery,
	}
	client.config = NewConfigApi(client)
	client.memoryCache = cache.NewMemoryCache(&cache.MemoryCacheConfig{
		Size:          1024,
		DefaultExpiry: 60,
	})
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
	enabled, err := c.client.memoryCache.Get(ctx, cacheKey)
	if err != nil {
		in := &proto.CheckConfigStatusRequest{
			ObjectName: objectName,
			DomainId:   domainId,
		}
		conn, err := c.client.getRandomClient()
		if err != nil {
			return false, err
		}
		res, err := conn.config.CheckConfigStatus(ctx, in)
		if err != nil {
			return false, err
		}
		c.client.memoryCache.Set(ctx, cacheKey, res.GetIsEnabled(), 120)
		if err != nil {
			return false, err
		}
		return res.GetIsEnabled(), nil
	}
	return enabled.Raw().(bool), nil
}

package client

import (
	"github.com/webitel/engine/discovery"
)

type Client struct {
	rabbit RabbitClient
	grpc   GrpcClient
}

func (c *Client) IsOpened() bool {
	return c.rabbit.IsOpened() && c.grpc.IsOpened()
}

// ! NewClient creates new client for logger.
// * rabbitUrl - connection string to rabbit1 server
// * clientId - name that will be recognized by consul
// * address - address to connect to consul server
func NewClient(rabbitUrl string, clientId string, consulAddress string) (*Client, error) {
	disc, err := discovery.NewServiceDiscovery(clientId, consulAddress, func() (bool, error) {
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	cli := &Client{grpc: NewGrpcClient(disc)}
	rab := NewRabbitClient(rabbitUrl, cli)
	cli.rabbit = rab
	return cli, nil
}

func (c *Client) Open() error {
	err := c.rabbit.Open()
	if err != nil {
		return err
	}
	err = c.grpc.Start()
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Close() {
	c.rabbit.Close()
	c.grpc.Stop()
}

func (c *Client) Rabbit() RabbitClient {
	return c.rabbit
}

func (c *Client) Grpc() GrpcClient {
	return c.grpc
}

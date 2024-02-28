package client

type Client struct {
	rabbit RabbitClient
	grpc   GrpcClient
}

func (c *Client) IsOpened() bool {
	return c.rabbit.IsOpened() && c.grpc.IsOpened()
}

// ! NewClient creates new client for logger.
// * rabbitUrl - connection string to rabbit1 server
// * consulAddress - address to connect to consul server
func NewClient(rabbitUrl string, consulAddress string) (*Client, error) {
	cli := &Client{grpc: NewGrpcClient(consulAddress)}
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

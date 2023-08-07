package client

import (
	"github.com/webitel/logger/proto"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Connection struct {
	name   string
	host   string
	client *grpc.ClientConn
	config proto.ConfigServiceClient
}

func NewConnection(name, url string) (*Connection, error) {
	var err error
	connection := &Connection{
		name: name,
		host: url,
	}

	connection.client, err = grpc.Dial(url, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(2*time.Second))

	if err != nil {
		return nil, err
	}

	connection.config = proto.NewConfigServiceClient(connection.client)

	return connection, nil
}

func (conn *Connection) Ready() bool {
	switch conn.client.GetState() {
	case connectivity.Idle, connectivity.Ready:
		return true
	}
	return false
}

func (conn *Connection) Name() string {
	return conn.name
}

func (conn *Connection) Close() error {
	err := conn.client.Close()
	if err != nil {
		return err
	}
	return nil
}

func (conn *Connection) Agent() proto.ConfigServiceClient {
	return conn.config
}

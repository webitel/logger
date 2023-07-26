package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Url string
}

type Action string

const (
	Create Action = "create"
	Update Action = "update"
	Delete Action = "delete"
	Read   Action = "read"
)

type RabbitClient interface {
	Open() error
	SendContext(ctx context.Context, domainId int, objectId int, userId int, userIp string, action Action, recordId int, newState []byte) error
	Close()
}

type rabbitClient struct {
	config  *Config
	conn    *amqp.Connection
	channel *amqp.Channel
	client  *Client
}

type Message struct {
	//ObjectId int    `json:"objectId,omitempty"`
	NewState []byte `json:"newState,omitempty"`
	UserId   int    `json:"userId,omitempty"`
	UserIp   string `json:"userIp,omitempty"`
	Action   string `json:"action,omitempty"`
	Date     int64  `json:"date,omitempty"`
	RecordId int    `json:"recordId,omitempty"`
	//DomainId int    `json:"domainId,omitempty"`
}

func NewRabbitClient(url string, client *Client) RabbitClient {
	return &rabbitClient{config: &Config{Url: url}, client: client}
}

func (c *rabbitClient) Open() error {
	conn, err := amqp.Dial(c.config.Url)
	if err != nil {
		return err
	}
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	err = channel.ExchangeDeclare(
		"logger", // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}
	queue, err := channel.QueueDeclare(
		"logger.service",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}
	err = channel.QueueBind(
		queue.Name, // queue name
		"logger.#", // routing key
		"logger",   // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}
	c.channel = channel
	c.conn = conn
	return nil
}

func (c *rabbitClient) Close() {
	c.channel.Close()
	c.conn.Close()
	c.channel = nil
	c.conn = nil
}

func (c *rabbitClient) SendContext(ctx context.Context, domainId int, objectId int, userId int, userIp string, action Action, recordId int, newState []byte) error {
	if c.channel == nil || c.conn == nil {
		return fmt.Errorf("connection not opened")
	}
	enabled, err := c.client.grpc.Config().CheckIsActive(context.Background(), domainId, objectId)
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}
	mess := &Message{UserId: userId, UserIp: userIp, Action: string(action), NewState: newState, Date: time.Now().Unix(), RecordId: recordId}
	result, err := json.Marshal(mess)
	if err != nil {
		return err
	}
	err = c.channel.PublishWithContext(
		ctx,
		"logger",
		fmt.Sprintf("logger.%d.%d", domainId, objectId),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        result,
		},
	)
	if err != nil {
		return err
	}
	return nil
}

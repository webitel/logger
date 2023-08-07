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
	CREATE_ACTION Action = "create"
	UPDATE_ACTION Action = "update"
	DELETE_ACTION Action = "delete"
	READ_ACTION   Action = "read"
)

type RabbitClient interface {
	Open() error
	SendContext(ctx context.Context, message *Message) error
	Close()
	IsOpened() bool
}

type rabbitClient struct {
	config   *Config
	conn     *amqp.Connection
	channel  *amqp.Channel
	client   *Client
	isOpened bool
}

type RequiredFields struct {
	UserId   int    `json:"userId,omitempty"`
	UserIp   string `json:"userIp,omitempty"`
	Action   string `json:"action,omitempty"`
	Date     int64  `json:"date,omitempty"`
	DomainId int64
	ObjectId int64
	//RecordId int    `json:"recordId,omitempty"`
}

type Message struct {
	RecordsStates  map[int][]byte `json:"records,omitempty"`
	NewState       []byte         `json:"newState,omitempty" `
	RecordId       int            `json:"recordId,omitempty"`
	RequiredFields `json:"requiredFields"`
	client         RabbitClient
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
	c.isOpened = true
	return nil
}

func (c *rabbitClient) Close() {
	c.channel.Close()
	c.conn.Close()
	c.channel = nil
	c.conn = nil
	c.isOpened = false
}

func (c *rabbitClient) IsOpened() bool {
	return c.isOpened
}

//func (c *rabbitClient) SendContext(ctx context.Context, domainId int, objectId int, userId int, userIp string, action Action, recordId int, newState []byte) error {
//	if c.channel == nil || c.conn == nil {
//		return fmt.Errorf("connection not opened")
//	}
//	enabled, err := c.client.grpc.Config().CheckIsActive(ctx, domainId, objectId)
//	if err != nil {
//		return err
//	}
//	if !enabled {
//		return nil
//	}
//
//	//mess := &Message{UserId: userId, UserIp: userIp, Action: string(action), NewState: newState, Date: time.Now().Unix(), RecordId: recordId}
//	//
//	//c.CreateAction(userId, userIp).Many()
//
//	result, err := json.Marshal(mess)
//	if err != nil {
//		return err
//	}
//	err = c.channel.PublishWithContext(
//		ctx,
//		"logger",
//		fmt.Sprintf("logger.%d.%d", domainId, objectId),
//		false,
//		false,
//		amqp.Publishing{
//			ContentType: "application/json",
//			Body:        result,
//		},
//	)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (c *rabbitClient) SendContext(ctx context.Context, message *Message) error {
	if c.IsOpened() {
		return fmt.Errorf("connection not opened")
	}
	enabled, err := c.client.Grpc().Config().CheckIsActive(ctx, int(message.RequiredFields.DomainId), int(message.RequiredFields.ObjectId))
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}

	//mess := &Message{UserId: userId, UserIp: userIp, Action: string(action), NewState: newState, Date: time.Now().Unix(), RecordId: recordId}
	//

	result, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = c.channel.PublishWithContext(
		ctx,
		"logger",
		fmt.Sprintf("logger.%d.%d", message.RequiredFields.DomainId, message.RequiredFields.ObjectId),
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

func (c *rabbitClient) CreateAction(domainId int64, objectId int64, userId int, userIp string) *Message {
	mess := &Message{RequiredFields: RequiredFields{
		UserId:   userId,
		UserIp:   userIp,
		Action:   string(CREATE_ACTION),
		Date:     time.Now().Unix(),
		DomainId: domainId,
		ObjectId: objectId,
	}, client: c}
	return mess
}

func (c *rabbitClient) UpdateAction(userId int, userIp string) *Message {
	mess := &Message{RequiredFields: RequiredFields{
		UserId: userId,
		UserIp: userIp,
		Action: string(UPDATE_ACTION),
		Date:   time.Now().Unix(),
	}, client: c}
	return mess
}

func (c *rabbitClient) DeleteAction(userId int, userIp string) *Message {
	mess := &Message{RequiredFields: RequiredFields{
		UserId: userId,
		UserIp: userIp,
		Action: string(DELETE_ACTION),
		Date:   time.Now().Unix(),
	}, client: c}
	return mess
}

func (c *Message) Many(recordsId []int, newStates [][]byte) *Message {
	m := make(map[int][]byte)
	if len(recordsId) != len(newStates) {
		return c
	}
	for i, v := range recordsId {
		m[v] = newStates[i]
	}
	c.RecordsStates = m
	return c
}

func (c *Message) One(recordId int, newState []byte) *Message {
	c.RecordId = recordId
	c.NewState = newState
	return c
}

func (c *Message) SendContext(ctx context.Context) error {
	err := c.client.SendContext(ctx, c)
	if err != nil {
		return err
	}
	return nil
}

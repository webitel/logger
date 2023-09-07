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

func (a Action) String() string {
	return string(a)
}

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
	CreateAction(domainId int64, objectName string, userId int, userIp string) *Message
	UpdateAction(domainId int64, objectName string, userId int, userIp string) *Message
	DeleteAction(domainId int64, objectName string, userId int, userIp string) *Message
}

type rabbitClient struct {
	config   *Config
	conn     *amqp.Connection
	channel  *amqp.Channel
	client   *Client
	isOpened bool
}

type RequiredFields struct {
	UserId     int    `json:"userId,omitempty"`
	UserIp     string `json:"userIp,omitempty"`
	Action     string `json:"action,omitempty"`
	Date       int64  `json:"date,omitempty"`
	DomainId   int64
	ObjectName string
	//RecordId int    `json:"recordId,omitempty"`
}

type Record struct {
	Id       int64  `json:"id,omitempty"`
	NewState []byte `json:"newState,omitempty"`
}

type Message struct {
	Records        []*Record `json:"records,omitempty"`
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
	_, err = channel.QueueDeclare(
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
	c.channel = channel
	c.conn = conn
	c.isOpened = true
	return nil
}

func (c *rabbitClient) Close() {
	c.conn.Close()
	c.channel = nil
	c.conn = nil
	c.isOpened = false
}

func (c *rabbitClient) IsOpened() bool {
	return c.isOpened
}

func (c *rabbitClient) SendContext(ctx context.Context, message *Message) error {
	if !c.IsOpened() {
		return fmt.Errorf("connection not opened")
	}
	enabled, err := c.client.Grpc().Config().CheckIsActive(ctx, message.RequiredFields.DomainId, message.RequiredFields.ObjectName)
	if err != nil {
		return err
	}
	if !enabled {
		return nil
	}

	if err := message.checkRecordsValidity(); err != nil {
		return err
	}

	result, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = c.channel.PublishWithContext(
		ctx,
		"logger",
		fmt.Sprintf("logger.%d.%s", message.RequiredFields.DomainId, message.RequiredFields.ObjectName),
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

func (c *rabbitClient) CreateAction(domainId int64, objectName string, userId int, userIp string) *Message {
	mess := &Message{RequiredFields: RequiredFields{
		UserId:     userId,
		UserIp:     userIp,
		Action:     string(CREATE_ACTION),
		Date:       time.Now().Unix(),
		DomainId:   domainId,
		ObjectName: objectName,
	}, client: c}
	return mess
}

func (c *rabbitClient) UpdateAction(domainId int64, objectName string, userId int, userIp string) *Message {
	mess := &Message{RequiredFields: RequiredFields{
		UserId:     userId,
		UserIp:     userIp,
		Action:     string(UPDATE_ACTION),
		Date:       time.Now().Unix(),
		DomainId:   domainId,
		ObjectName: objectName,
	}, client: c}
	return mess
}

func (c *rabbitClient) DeleteAction(domainId int64, objectName string, userId int, userIp string) *Message {
	mess := &Message{RequiredFields: RequiredFields{
		UserId:     userId,
		UserIp:     userIp,
		Action:     string(DELETE_ACTION),
		Date:       time.Now().Unix(),
		DomainId:   domainId,
		ObjectName: objectName,
	}, client: c}
	return mess
}

func (c *Message) Many(records []*Record) *Message {
	if len(records) == 0 {
		return c
	}
	c.Records = records
	return c
}

func (c *Message) checkRecordsValidity() error {
	if c.Records == nil {
		return fmt.Errorf("logger: no records data in message")
	}
	var canNil bool
	switch c.Action {
	case CREATE_ACTION.String(), UPDATE_ACTION.String():
		canNil = false
	case DELETE_ACTION.String():
		canNil = true
	}
	if !canNil {
		for _, record := range c.Records {
			if record.NewState == nil || len(record.NewState) == 0 {
				return fmt.Errorf("logger: record has no data ( id: %s )", record.Id)
			}
		}
	}

	return nil
}

func (c *Message) One(record *Record) *Message {
	if record == nil {
		return c
	}
	c.Records = append(c.Records, record)
	return c
}

func (c *Message) SendContext(ctx context.Context) error {
	//if err := c.checkRecordsValidity(); err != nil {
	//	return err
	//}
	err := c.client.SendContext(ctx, c)
	if err != nil {
		return err
	}
	return nil
}

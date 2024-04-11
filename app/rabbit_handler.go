package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/webitel/logger/model"
	"github.com/webitel/wlog"

	amqp "github.com/rabbitmq/amqp091-go"
	errors "github.com/webitel/engine/model"
)

type Handler struct {
	app *App
}
type Message struct {
	Records        []*Record `json:"records,omitempty"`
	RequiredFields `json:"requiredFields"`
}

type RequiredFields struct {
	UserId int    `json:"userId,omitempty"`
	UserIp string `json:"userIp,omitempty"`
	Action string `json:"action,omitempty"`
	Date   int64  `json:"date,omitempty"`
}

type Record struct {
	Id       int64     `json:"id,omitempty"`
	NewState BytesJSON `json:"newState,omitempty"`
}

type BytesJSON struct {
	Body []byte
}

func (b *BytesJSON) GetBody() []byte {
	return b.Body
}

func (b *BytesJSON) UnmarshalJSON(input []byte) error {
	b.Body = input
	return nil
}

func NewHandler(app *App) (*Handler, errors.AppError) {
	if app == nil {
		return nil, errors.NewInternalError("rabbit.handler.new_handler.arguments_check.app_nil", "can't configure handler, app is nil")
	}
	return &Handler{app: app}, nil
}

func (h *Handler) Handle(ctx context.Context, message *amqp.Delivery) errors.AppError {
	var (
		m      Message
		domain int64
		object string
	)
	err := json.Unmarshal(message.Body, &m)
	if err != nil {
		wlog.Debug(fmt.Sprintf("error unmarshalling message. details: %s", err.Error()))
		return nil
		//return errors.NewInternalError("rabbit.handler.handle.json_unmarshal.error", err.Error())
	}

	splittedKey := strings.Split(message.RoutingKey, ".")
	if len(splittedKey) >= 3 {
		domain, _ = strconv.ParseInt(splittedKey[1], 10, 64)
		object = splittedKey[2]
	}
	if m.Records != nil {
		var rabbitMessages []*model.RabbitMessage
		for _, v := range m.Records {
			rabbitMessage := &model.RabbitMessage{
				//ObjectId: object,
				NewState: v.NewState.GetBody(),
				UserId:   m.UserId,
				UserIp:   m.UserIp,
				Action:   m.Action,
				Date:     m.Date,
				//DomainId: domain,
				RecordId: v.Id,
				Schema:   object,
			}
			rabbitMessages = append(rabbitMessages, rabbitMessage)
		}
		appErr := h.app.InsertLogByRabbitMessageBulk(ctx, rabbitMessages, domain, object)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

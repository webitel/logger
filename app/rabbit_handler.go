package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/webitel/logger/model"
	"github.com/webitel/logger/pkg/client"
	"github.com/webitel/wlog"

	amqp "github.com/rabbitmq/amqp091-go"
	errors "github.com/webitel/engine/model"
)

type Handler struct {
	app *App
}

func NewHandler(app *App) (*Handler, errors.AppError) {
	if app == nil {
		return nil, errors.NewInternalError("rabbit1.handler.new_handler.arguments_check.app_nil", "can't configure handler, app is nil")
	}
	return &Handler{app: app}, nil
}

func (h *Handler) Handle(ctx context.Context, message *amqp.Delivery) errors.AppError {
	var (
		m      client.Message
		domain int
		object int
	)
	err := json.Unmarshal(message.Body, &m)
	if err != nil {
		wlog.Debug(fmt.Sprintf("error unmarshalling message. details: %s", err.Error()))
		return nil
		//return errors.NewInternalError("rabbit1.handler.handle.json_unmarshal.error", err.Error())
	}

	splittedKey := strings.Split(message.RoutingKey, ".")
	if len(splittedKey) >= 3 {
		domain, _ = strconv.Atoi(splittedKey[1])
		object, _ = strconv.Atoi(splittedKey[2])
	}
	if m.RecordsStates != nil {
		var rabbitMessages []*model.RabbitMessage
		for i, v := range m.RecordsStates {
			rabbitMessage := &model.RabbitMessage{
				//ObjectId: object,
				NewState: v,
				UserId:   m.UserId,
				UserIp:   m.UserIp,
				Action:   m.Action,
				Date:     m.Date,
				//DomainId: domain,
				RecordId: i,
			}
			rabbitMessages = append(rabbitMessages, rabbitMessage)
		}
		appErr := h.InsertLogByRabbitMessageBulk(ctx, rabbitMessages, domain, object)
		if appErr != nil {
			return appErr
		}
	} else {
		rabbitMessage := &model.RabbitMessage{
			//	ObjectId: domain,
			NewState: m.NewState,
			UserId:   m.UserId,
			UserIp:   m.UserIp,
			Action:   m.Action,
			Date:     m.Date,
			//DomainId: domain,
			RecordId: m.RecordId,
		}
		appErr := h.InsertLogByRabbitMessage(ctx, rabbitMessage, domain, object)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

package app

import (
	"context"
	"encoding/json"
	"fmt"
	proto "github.com/webitel/protos/logger"
	"strconv"
	"strings"
	"time"

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
		return nil, errors.NewInternalError("rabbit.handler.new_handler.arguments_check.app_nil", "can't configure handler, app is nil")
	}
	return &Handler{app: app}, nil
}

func (h *Handler) Handle(ctx context.Context, message *amqp.Delivery) errors.AppError {
	var (
		m      client.Message
		domain int64
		object string
	)
	err := json.Unmarshal(message.Body, &m)
	if err != nil {
		wlog.Debug(fmt.Sprintf("rabbit_handler: error unmarshalling message. details: %s", err.Error()))
		return nil
		//return errors.NewInternalError("rabbit.handler.handle.json_unmarshal.error", err.Error())
	}

	splittedKey := strings.Split(message.RoutingKey, ".")
	if len(splittedKey) >= 3 {
		domain, _ = strconv.ParseInt(splittedKey[1], 10, 64)
		object = splittedKey[2]
	}
	if m.Records != nil {
		var (
			rabbitMessages []*model.RabbitMessage
		)

		for _, v := range m.Records {

			rabbitMessage := &model.RabbitMessage{
				//ObjectId: object,
				NewState: v.NewState,
				UserId:   m.UserId,
				UserIp:   m.UserIp,
				Action:   m.Action,
				Date:     m.Date,
				//DomainId: domain,
				RecordId: v.Id,
				Schema:   object,
			}
			rabbitMessages = append(rabbitMessages, rabbitMessage)
			if object == proto.AvailableSystemObjects_schema.String() {
				userId, err := model.NewNullInt(m.UserId)
				if err != nil {
					return errors.NewBadRequestError("rabbit.handler.new_handler.arguments_conv.error", err.Error())
				}
				appErr := h.app.storage.SchemaVersion().Insert(ctx, &model.SchemaVersion{
					SchemaId:   v.Id,
					CreatedOn:  time.Unix(m.Date, 0),
					CreatedBy:  model.Lookup{Id: userId},
					ObjectData: v.NewState,
					Note:       model.NewNullString(v.Note),
				})
				if appErr != nil {
					return appErr
				}
			}
		}

		status, appErr := h.app.CheckConfigStatus(ctx, &proto.CheckConfigStatusRequest{
			ObjectName: object,
			DomainId:   domain,
		})
		if appErr != nil {
			wlog.Debug(fmt.Sprintf("rabbit_handler: error getting config state. details: %s", appErr.Error()))
			return nil
		}
		if status.IsEnabled {
			appErr := h.app.InsertLogByRabbitMessageBulk(ctx, rabbitMessages, domain, object)
			if appErr != nil {
				return appErr
			}
		}

	}

	return nil
}

package rabbit

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"webitel_logger/app"
	"webitel_logger/model"

	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	errors "github.com/webitel/engine/model"
)

type Handler struct {
	app *app.App
}

func NewHandler(app *app.App) (*Handler, errors.AppError) {
	if app == nil {
		return nil, errors.NewInternalError("rabbit.handler.new_handler.arguments_check.app_nil", "can't configure handler, app is nil")
	}
	return &Handler{app: app}, nil
}

func (h *Handler) Handle(ctx context.Context, message *amqp.Delivery) errors.AppError {
	var model model.Message
	err := json.Unmarshal(message.Body, &model)
	if err != nil {
		log.Printf("error unmarshalling message. details: %s", err.Error())
		return nil
		//return errors.NewInternalError("rabbit.handler.handle.json_unmarshal.error", err.Error())
	}
	splittedKey := strings.Split(message.RoutingKey, ".")
	if len(splittedKey) >= 3 {
		model.DomainId, _ = strconv.Atoi(splittedKey[1])
		model.ObjectId, _ = strconv.Atoi(splittedKey[2])
	}
	return h.app.InsertLogByRabbitMessage(ctx, &model)
}

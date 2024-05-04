package controller

import (
	"context"
	"golang-note/helper"
	"golang-note/middleware"
	modelresponse "golang-note/model/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rabbitmq/amqp091-go"
)

type RabbitmqController interface {
	PushMessage(c echo.Context) error
}

type RabbitmqControllerImplementation struct {
	Channel *amqp091.Channel
}

func NewRabbitmqController(channel *amqp091.Channel) RabbitmqController {
	return &RabbitmqControllerImplementation{
		Channel: channel,
	}
}

func (controller *RabbitmqControllerImplementation) PushMessage(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	key := c.Param("key")
	messageParam := c.Param("message")
	message := amqp091.Publishing{
		Headers: amqp091.Table{
			"sample": "value",
		},
		Body: []byte("message: " + messageParam),
	}
	ctx := context.Background()
	err := controller.Channel.PublishWithContext(ctx, "notification", key, false, false, message)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
	}

	return modelresponse.ToResponse(c, http.StatusOK, requestId, "successfully publish message", "")
}

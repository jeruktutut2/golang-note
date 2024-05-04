package controller

import (
	"golang-note/exception"
	"golang-note/helper"
	"golang-note/middleware"
	modelresponse "golang-note/model/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lovoo/goka"
)

type KafkaController interface {
	EmitEmailMessage(c echo.Context) error
	EmitTextMessage(c echo.Context) error
}

type KafkaControllerImplementation struct {
	EmailEmitter       *goka.Emitter
	TextMessageEmitter *goka.Emitter
}

func NewKafkaController(emailEmmiter *goka.Emitter, textMessageEmitter *goka.Emitter) KafkaController {
	return &KafkaControllerImplementation{
		EmailEmitter:       emailEmmiter,
		TextMessageEmitter: textMessageEmitter,
	}
}

func (controller *KafkaControllerImplementation) EmitEmailMessage(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	message := c.Param("message")
	err := controller.EmailEmitter.EmitSync("email", message)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewInternalServerErrorException()
		return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
	}
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "successfully sent message", "")
}

func (controller *KafkaControllerImplementation) EmitTextMessage(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	message := c.Param("message")
	err := controller.TextMessageEmitter.EmitSync("text-message", message)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
	}
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "successfully sent message", "")
}

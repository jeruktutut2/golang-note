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
	EmitMessage(c echo.Context) error
}

type KafkaControllerImplementation struct {
	Emmiter *goka.Emitter
}

func NewKafkaController(emmiter *goka.Emitter) KafkaController {
	return &KafkaControllerImplementation{
		Emmiter: emmiter,
	}
}

func (controller *KafkaControllerImplementation) EmitMessage(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	key := c.Param("key")
	message := c.Param("message")
	err := controller.Emmiter.EmitSync(key, message)
	if err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewInternalServerErrorException()
		return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
	}
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "successfully sent message", "")
}

package controller

import (
	"context"
	"golang-note/exception"
	"golang-note/middleware"
	modelresponse "golang-note/model/response"
	"golang-note/service"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ContextController interface {
	Timeout(c echo.Context) error
	CustomTimeout(c echo.Context) error
	TimeoutDb(c echo.Context) error
	TimeoutTx(c echo.Context) error
}

type ContextControllerImplementation struct {
	ContextService service.ContextService
}

func NewContextController(contextService service.ContextService) ContextController {
	return &ContextControllerImplementation{
		ContextService: contextService,
	}
}

func (controller *ContextControllerImplementation) Timeout(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	str, err := controller.ContextService.Timeout(c.Request().Context(), requestId)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}
	return modelresponse.ToResponse(c, http.StatusOK, requestId, str, "")
}

func (controller *ContextControllerImplementation) CustomTimeout(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	ctx, cancel := context.WithTimeout(c.Request().Context(), time.Duration(3)*time.Second)
	defer cancel()

	c.SetRequest(c.Request().WithContext(ctx))
	str, err := controller.ContextService.Timeout(c.Request().Context(), requestId)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}
	return modelresponse.ToResponse(c, http.StatusOK, requestId, str, "")
}

func (controller *ContextControllerImplementation) TimeoutDb(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	str, err := controller.ContextService.CheckContext(c.Request().Context(), requestId)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}
	return modelresponse.ToResponse(c, http.StatusOK, requestId, str, "")
}

func (controller *ContextControllerImplementation) TimeoutTx(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	str, err := controller.ContextService.CheckContextTx(c.Request().Context(), requestId)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}
	return modelresponse.ToResponse(c, http.StatusOK, requestId, str, "")
}

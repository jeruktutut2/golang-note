package controller

import (
	"golang-note/middleware"
	modelresponse "golang-note/model/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RestController interface {
	Get(c echo.Context) error
	Post(c echo.Context) error
	Put(c echo.Context) error
	Patch(c echo.Context) error
	Delete(c echo.Context) error
	Head(c echo.Context) error
	Options(c echo.Context) error
}

type RestControllerImplementation struct {
}

func NewRestController() RestController {
	return &RestControllerImplementation{}
}

func (controller *RestControllerImplementation) Get(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "Get", "")
}

func (controller *RestControllerImplementation) Post(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "Post", "")
}

func (controller *RestControllerImplementation) Put(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "Put", "")
}

func (controller *RestControllerImplementation) Patch(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "Patch", "")
}

func (controller *RestControllerImplementation) Delete(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "Delete", "")
}

func (controller *RestControllerImplementation) Head(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "Head", "")
}

func (controllerImplementation *RestControllerImplementation) Options(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "Options", "")
}

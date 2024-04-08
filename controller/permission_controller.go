package controller

import (
	"golang-note/middleware"
	modelresponse "golang-note/model/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PermissionController interface {
	Permission(c echo.Context) error
	NoPermission(c echo.Context) error
	PermissionRedis(c echo.Context) error
	NoPermissionRedis(c echo.Context) error
}

type PermissionControllerImplementation struct {
}

func NewPermissionController() PermissionController {
	return &PermissionControllerImplementation{}
}

func (controller *PermissionControllerImplementation) Permission(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	permissions := c.Request().Context().Value(middleware.PermissionKey).([]interface{})
	return modelresponse.ToResponse(c, http.StatusOK, requestId, permissions, "")
}

func (controller *PermissionControllerImplementation) NoPermission(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	permissions := c.Request().Context().Value(middleware.PermissionKey).([]interface{})
	return modelresponse.ToResponse(c, http.StatusOK, requestId, permissions, "")
}

func (controller *PermissionControllerImplementation) PermissionRedis(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "permission", "")
}

func (controller *PermissionControllerImplementation) NoPermissionRedis(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "permission", "")
}

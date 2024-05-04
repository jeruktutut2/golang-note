package controller

import (
	"golang-note/exception"
	"golang-note/helper"
	"golang-note/middleware"
	modelrequest "golang-note/model/request"
	modelresponse "golang-note/model/response"
	"golang-note/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UserController interface {
	Login(c echo.Context) error
	RefreshToken(c echo.Context) error
	LoginRedis(c echo.Context) error
	CheckLoginRedis(c echo.Context) error
	LoginMap(c echo.Context) error
	CheckLoginMap(c echo.Context) error
	Logout(c echo.Context) error
	LogoutRedis(c echo.Context) error
	LogoutMap(c echo.Context) error
}

type UserControllerImplementaion struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImplementaion{
		UserService: userService,
	}
}

func (controller *UserControllerImplementaion) Login(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	loginUserRequest := modelrequest.LoginRequest{}
	if err := c.Bind(&loginUserRequest); err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewBadRequestException(err.Error())
		return exception.ErrorHandler(c, requestId, err)
	}
	loginResponse, accessToken, refreshToken, err := controller.UserService.Login(c.Request().Context(), requestId, loginUserRequest)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = "Bearer " + accessToken
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	cookie = new(http.Cookie)
	cookie.Name = "X-REFRESH-TOKEN"
	cookie.Value = refreshToken
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return modelresponse.ToResponse(c, http.StatusOK, requestId, loginResponse, "")
}

func (controller *UserControllerImplementaion) RefreshToken(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	username := c.Request().Context().Value(middleware.UsernameKey).(string)
	xRefreshToken := c.Request().Context().Value(middleware.XRefreshTokenKey).(string)

	accessToken, err := controller.UserService.RefreshToken(c.Request().Context(), requestId, username, xRefreshToken)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = "Bearer " + accessToken
	cookie.Path = "/"
	cookie.Domain = "localhost"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return modelresponse.ToResponse(c, http.StatusOK, requestId, "successfully set new access token", "")
}

func (controller *UserControllerImplementaion) LoginRedis(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	loginRequest := modelrequest.LoginRequest{}
	if err := c.Bind(&loginRequest); err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewBadRequestException(err.Error())
		return exception.ErrorHandler(c, requestId, err)
	}
	loginResponse, token, err := controller.UserService.LoginRedis(c.Request().Context(), requestId, loginRequest)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = token
	cookie.Path = "/"
	// cookie.Domain = "localhost"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return modelresponse.ToResponse(c, http.StatusOK, requestId, loginResponse, "")
}

func (controller *UserControllerImplementaion) CheckLoginRedis(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "check login redis", "")
}

func (controller *UserControllerImplementaion) LoginMap(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	loginRequest := modelrequest.LoginRequest{}
	if err := c.Bind(&loginRequest); err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewBadRequestException(err.Error())
		return exception.ErrorHandler(c, requestId, err)
	}
	loginResponse, token, err := controller.UserService.LoginMap(c.Request().Context(), requestId, loginRequest)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = token
	cookie.Path = "/"
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return modelresponse.ToResponse(c, http.StatusOK, requestId, loginResponse, "")
}

func (controller *UserControllerImplementaion) CheckLoginMap(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "check login map", "")
}

func (controller *UserControllerImplementaion) Logout(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	coockie := new(http.Cookie)
	coockie.Name = "Authorization"
	coockie.Value = ""
	coockie.Path = "/"
	coockie.MaxAge = -1
	coockie.HttpOnly = true
	c.SetCookie(coockie)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, "successfully logout", "")
}

func (controller *UserControllerImplementaion) LogoutRedis(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	key := c.Request().Context().Value(middleware.TokenIdKey).(string)
	logoutRedisResponse, err := controller.UserService.LogoutRedis(c.Request().Context(), requestId, key)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}
	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	cookie.HttpOnly = true
	c.SetCookie(cookie)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, logoutRedisResponse, "")
}

func (controller *UserControllerImplementaion) LogoutMap(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	sessionIs := c.Request().Context().Value(middleware.SessionKey).(string)
	logoutMapResponse, err := controller.UserService.LogoutMap(c.Request().Context(), sessionIs)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}
	return modelresponse.ToResponse(c, http.StatusOK, requestId, logoutMapResponse, "")
}

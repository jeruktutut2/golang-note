package middleware

import (
	"context"
	"errors"
	"golang-note/globals"
	"golang-note/helper"
	modelresponse "golang-note/model/response"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AuthenticateMap(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := c.Request().Context().Value(RequestIdKey).(string)
		cookie, err := c.Cookie("Authorization")
		if err != nil && err != http.ErrNoCookie {
			helper.PrintLogToTerminal(err, requestId)
			return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
		} else if err != nil && err == http.ErrNoCookie {
			helper.PrintLogToTerminal(err, requestId)
			return modelresponse.ToResponse(c, http.StatusUnauthorized, requestId, "", "unauthorized access")
		}

		session := cookie.Value
		if globals.Session[session] == "" {
			err = errors.New("cannot find cookie in redis")
			helper.PrintLogToTerminal(err, requestId)
			return modelresponse.ToResponse(c, http.StatusUnauthorized, requestId, "", "unauthorized access")
		}
		ctx := context.WithValue(c.Request().Context(), SessionKey, session)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}

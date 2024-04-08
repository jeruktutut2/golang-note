package middleware

import (
	"context"
	"errors"
	"golang-note/helper"
	modelresponse "golang-note/model/response"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func AuthenticateRedis(client *redis.Client) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
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
			key := cookie.Value
			_, err = client.Get(c.Request().Context(), key).Result()
			if err != nil && err != redis.Nil {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusUnauthorized, requestId, "", "unauthorized access")
			} else if err == redis.Nil {
				err = errors.New("cannot find cookie in redis")
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusUnauthorized, requestId, "", "unauthorized access")
			}
			ctx := context.WithValue(c.Request().Context(), TokenIdKey, key)
			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

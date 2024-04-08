package middleware

import (
	"golang-note/helper"
	modelresponse "golang-note/model/response"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func CheckPermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := c.Request().Context().Value(RequestIdKey).(string)
		permissions := c.Request().Context().Value(PermissionKey).([]interface{})
		var permissinoMatch bool
	permissionsLoop:
		for _, permission := range permissions {
			prmsn := permission.(string)
			if prmsn == "CREATE_BOOK" {
				permissinoMatch = true
				break permissionsLoop
				// next(writer, request, params)
			}
		}
		if permissinoMatch {
			return next(c)
		} else {
			return modelresponse.ToResponse(c, http.StatusForbidden, requestId, "", "not permitted")
		}
	}
}

func Nopermission(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		requestId := c.Request().Context().Value(RequestIdKey).(string)
		permissions := c.Request().Context().Value(PermissionKey).([]interface{})
		var permissinoMatch bool
	permissionsLoop:
		for _, permission := range permissions {
			prmsn := permission.(string)
			if prmsn == "nopermission" {
				permissinoMatch = true
				break permissionsLoop
			}
		}
		if permissinoMatch {
			return next(c)
		} else {
			return modelresponse.ToResponse(c, http.StatusForbidden, requestId, "", "not permitted")
		}
	}
}

func CheckPermissionRedis(client *redis.Client) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestId := c.Request().Context().Value(RequestIdKey).(string)
			tokenId := c.Request().Context().Value(TokenIdKey).(string)
			permissionsString, err := client.Get(c.Request().Context(), tokenId).Result()
			if err != nil && err != redis.Nil {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
			} else if err == redis.Nil {
				return modelresponse.ToResponse(c, http.StatusForbidden, requestId, "", "not permitted")
			}
			permissions := strings.Split(permissionsString, ",")
			var permissinoMatch bool
		permissionsLoop:
			for _, permission := range permissions {
				if permission == "CREATE_BOOK" {
					permissinoMatch = true
					break permissionsLoop
				}
			}
			if permissinoMatch {
				return next(c)
			} else {
				return modelresponse.ToResponse(c, http.StatusForbidden, requestId, "", "not permitted")
			}
		}
	}
}

func CheckNoPermossionRedis(client *redis.Client) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestId := c.Request().Context().Value(RequestIdKey).(string)
			tokenId := c.Request().Context().Value(TokenIdKey).(string)
			permissionsString, err := client.Get(c.Request().Context(), tokenId).Result()
			if err != nil && err != redis.Nil {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
			} else if err == redis.Nil {
				return modelresponse.ToResponse(c, http.StatusForbidden, requestId, "", "not permitted")
			}
			permissions := strings.Split(permissionsString, ",")
			var permissinoMatch bool
		permissionsLoop:
			for _, permission := range permissions {
				if permission == "nopermission" {
					permissinoMatch = true
					break permissionsLoop
				}
			}
			if permissinoMatch {
				return next(c)
			} else {
				return modelresponse.ToResponse(c, http.StatusForbidden, requestId, "", "not permitted")
			}
		}
	}
}

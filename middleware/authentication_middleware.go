package middleware

import (
	"context"
	"errors"
	"fmt"
	"golang-note/helper"
	modelresponse "golang-note/model/response"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func Authenticate(jwtKey string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			requestId := c.Request().Context().Value(RequestIdKey).(string)
			authorizationCockie, err := c.Cookie("Authorization")
			if err != nil {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
			}
			authorizationToken := authorizationCockie.Value
			authorizationTokenSplit := strings.Split(authorizationToken, " ")
			tokenJwt := authorizationTokenSplit[1]

			token, err := jwt.Parse(tokenJwt, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtKey), nil
			})
			if errors.Is(err, jwt.ErrTokenMalformed) {
				// fmt.Println("That's not even a token")
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusUnauthorized, requestId, "", "token yang diberikan tidak valid atau bukan merupakan token JWT")
			} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
				// fmt.Println("Invalid signature")
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusForbidden, requestId, "", "signature token JWT tidak valid atau gagal")
			} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
				// fmt.Println("Timing is everything")
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusUnauthorized, requestId, "", "token yang diberikan sudah expire")
			} else if err != nil {
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				ctx := context.WithValue(c.Request().Context(), UsernameKey, claims["username"])
				c.SetRequest(c.Request().WithContext(ctx))
				ctx = context.WithValue(c.Request().Context(), PermissionKey, claims["permissions"])
				c.SetRequest(c.Request().WithContext(ctx))

				xRefreshToken, err := c.Cookie("X-REFRESH-TOKEN")
				if err != nil {
					helper.PrintLogToTerminal(err, requestId)
					return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "internal server error")
				}
				ctx = context.WithValue(c.Request().Context(), XRefreshTokenKey, xRefreshToken.Value)
				c.SetRequest(c.Request().WithContext(ctx))

				return next(c)
			} else {
				err = errors.New("token claims is not provided")
				helper.PrintLogToTerminal(err, requestId)
				return modelresponse.ToResponse(c, http.StatusInternalServerError, requestId, "", "token claims is not provided")
			}

			// next(writer, request)
		}
	}
}

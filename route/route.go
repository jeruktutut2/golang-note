package route

import (
	"golang-note/controller"
	"golang-note/middleware"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

func RestRoute(e *echo.Echo, controller controller.RestController) {
	e.GET("/api/v1/get/:nama", controller.Get)
	e.POST("/api/v1/get/:nama", controller.Post)
	e.PUT("/api/v1/get/:nama", controller.Put)
	e.PATCH("/api/v1/get/:nama", controller.Patch)
	e.DELETE("/api/v1/get/:nama", controller.Delete)
	e.HEAD("/api/v1/get/:nama", controller.Head)
	e.OPTIONS("/api/v1/get/:nama", controller.Options)
}

func UserRoute(e *echo.Echo, jwtKey string, client *redis.Client, controller controller.UserController) {
	e.POST("/api/v1/login", controller.Login)
	e.POST("/api/v1/refresh-token", controller.RefreshToken, middleware.Authenticate(jwtKey))
	e.POST("/api/v1/login-redis", controller.LoginRedis)
	e.GET("/api/v1/check-login-redis", controller.CheckLoginRedis, middleware.AuthenticateRedis(client))
	e.POST("/api/v1/login-map", controller.LoginMap)
	e.GET("/api/v1/check-login-map", controller.CheckLoginMap, middleware.AuthenticateMap)
}

func ContextRoute(e *echo.Echo, controller controller.ContextController) {
	e.GET("/api/v1/timeout", controller.Timeout)
	e.GET("/api/v1/custom-timeout", controller.CustomTimeout)
	e.GET("/api/v1/timeout-db", controller.TimeoutDb)
	e.GET("/api/v1/timeout-tx", controller.TimeoutTx)
}

func PermissionRoute(e *echo.Echo, jwtKey string, client *redis.Client, controller controller.PermissionController) {
	e.GET("/api/v1/permission", controller.Permission, middleware.CheckPermission, middleware.Authenticate(jwtKey))
	e.GET("/api/v1/nopermission", controller.NoPermission, middleware.Nopermission, middleware.Authenticate(jwtKey))
	e.GET("/api/v1/permission-redis", controller.PermissionRedis, middleware.CheckPermissionRedis(client), middleware.AuthenticateRedis(client))
	e.GET("/api/v1/nopermission-redis", controller.NoPermissionRedis, middleware.CheckNoPermossionRedis(client), middleware.AuthenticateRedis(client))
}

func PdfRoute(e *echo.Echo, controller controller.PdfController) {
	e.GET("/api/v1/pdf", controller.GetPdf)
}

func KafkaRoute(e *echo.Echo, controller controller.KafkaController) {
	e.GET("/api/v1/kafka/:key/:message", controller.EmitMessage)
}

func StreamJsonRoute(e *echo.Echo, controller controller.StreamJsonController) {
	e.GET("/api/v1/stream-json", controller.Stream)
}

func ServerSentEventRoute(e *echo.Echo, controller controller.ServerSentEventController) {
	e.GET("/api/v1/handle-sse", controller.HandleSSE)
	e.GET("/api/v1/send-message/:message", controller.SendMessage)
}

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
	e.POST("/api/v1/user/login", controller.Login)
	e.POST("/api/v1/user/refresh-token", controller.RefreshToken, middleware.Authenticate(jwtKey))
	e.POST("/api/v1/user/login-redis", controller.LoginRedis)
	e.GET("/api/v1/user/check-login-redis", controller.CheckLoginRedis, middleware.AuthenticateRedis(client))
	e.POST("/api/v1/user/login-map", controller.LoginMap)
	e.GET("/api/v1/user/check-login-map", controller.CheckLoginMap, middleware.AuthenticateMap)
	e.POST("/api/v1/user/logout", controller.Logout)
	e.POST("/api/v1/user/logout-redis", controller.LogoutRedis, middleware.AuthenticateRedis(client))
	e.POST("/api/v1/user/logout-map", controller.LogoutMap, middleware.AuthenticateMap)
}

func ContextRoute(e *echo.Echo, controller controller.ContextController) {
	e.GET("/api/v1/context/timeout", controller.Timeout)
	e.GET("/api/v1/context/custom-timeout", controller.CustomTimeout)
	e.GET("/api/v1/context/timeout-db", controller.TimeoutDb)
	e.GET("/api/v1/context/timeout-tx", controller.TimeoutTx)
}

func PermissionRoute(e *echo.Echo, jwtKey string, client *redis.Client, controller controller.PermissionController) {
	e.GET("/api/v1/permission/check-permission", controller.Permission, middleware.Authenticate(jwtKey), middleware.CheckPermission)
	e.GET("/api/v1/permission/check-nopermission", controller.NoPermission, middleware.Authenticate(jwtKey), middleware.Nopermission)
	e.GET("/api/v1/permission/check-permission-redis", controller.PermissionRedis, middleware.AuthenticateRedis(client), middleware.CheckPermissionRedis(client))
	e.GET("/api/v1/permission/check-nopermission-redis", controller.NoPermissionRedis, middleware.AuthenticateRedis(client), middleware.CheckNoPermossionRedis(client))
}

func PdfRoute(e *echo.Echo, controller controller.PdfController) {
	e.GET("/api/v1/pdf", controller.GetPdf)
}

func KafkaRoute(e *echo.Echo, controller controller.KafkaController) {
	e.GET("/api/v1/kafka/email/:message", controller.EmitEmailMessage)
	e.GET("/api/v1/kafka/text-message/:message", controller.EmitTextMessage)
}

func StreamJsonRoute(e *echo.Echo, controller controller.StreamJsonController) {
	e.GET("/api/v1/stream-json", controller.Stream)
	e.GET("/api/v1/stream-without-channel", controller.StreamWithoutChannel)
	e.GET("/api/v1/stream-with-sse", controller.StreamWithSSE)
}

func ServerSentEventRoute(e *echo.Echo, controller controller.ServerSentEventController) {
	e.GET("/api/v1/sse/handle-sse", controller.HandleSSE)
	e.GET("/api/v1/sse/send-message/:message", controller.SendMessage)
	e.GET("/api/v1/sse/handle-sse-without-channel", controller.HandleSSEWithoutChannel)
	e.GET("/api/v1/sse/send-message-without-channel/:message", controller.SendMessageWithoutChannel)
}

func RabbitmqController(e *echo.Echo, controller controller.RabbitmqController) {
	e.GET("/api/v1/rabbitmq/push-message/:key/:message", controller.PushMessage)
}

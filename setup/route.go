package setup

import (
	"golang-note/configuration"
	"golang-note/route"
	"golang-note/util"

	"github.com/labstack/echo/v4"
)

func Route(e *echo.Echo, config *configuration.Configuration, redis util.RedisUtil, controller ControllerSetup) {
	route.UserRoute(e, config.JwtKey, redis.GetClient(), controller.UserController)
	route.RestRoute(e, controller.RestController)
	route.ContextRoute(e, controller.ContextController)
	route.PermissionRoute(e, config.JwtKey, redis.GetClient(), controller.PermissionController)
	route.PdfRoute(e, controller.PdfController)
	route.KafkaRoute(e, controller.KafkaController)
	route.StreamJsonRoute(e, controller.StreamJsonController)
	route.ServerSentEventRoute(e, controller.ServerSentEventController)
	route.RabbitmqController(e, controller.RabbitmqController)
}

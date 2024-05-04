package setup

import (
	"golang-note/controller"

	"github.com/lovoo/goka"
	"github.com/rabbitmq/amqp091-go"
)

type ControllerSetup struct {
	UserController            controller.UserController
	RestController            controller.RestController
	ContextController         controller.ContextController
	PermissionController      controller.PermissionController
	PdfController             controller.PdfController
	KafkaController           controller.KafkaController
	StreamJsonController      controller.StreamJsonController
	ServerSentEventController controller.ServerSentEventController
	RabbitmqController        controller.RabbitmqController
}

func Controller(emailEmitter *goka.Emitter, textMessageEmitter *goka.Emitter, channel *amqp091.Channel, service ServiceSetup) ControllerSetup {
	return ControllerSetup{
		UserController:            controller.NewUserController(service.UserService),
		RestController:            controller.NewRestController(),
		ContextController:         controller.NewContextController(service.ContextService),
		PermissionController:      controller.NewPermissionController(),
		PdfController:             controller.NewPdfController(service.PdfService),
		KafkaController:           controller.NewKafkaController(emailEmitter, textMessageEmitter),
		StreamJsonController:      controller.NewStreamJsonController(service.StreamJsonService),
		ServerSentEventController: controller.NewServerSentEventController(),
		RabbitmqController:        controller.NewRabbitmqController(channel),
	}
}

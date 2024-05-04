package util

import (
	"golang-note/configuration"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func NewRabbitmqConneciton(config *configuration.Configuration) (connection *amqp091.Connection, channel *amqp091.Channel) {
	connection, err := amqp091.Dial("amqp://" + config.RabbitmqUsername + ":" + config.RabbitmqPassword + "@" + config.RabbitmqHost + ":" + config.RabbitmqPort + "/")
	if err != nil {
		log.Fatalln("error when connecting to rabbitmq server:", err)
	}
	// defer connection.Close()
	channel, err = connection.Channel()
	if err != nil {
		log.Fatalln("error when creating channel:", err)
	}
	return
}

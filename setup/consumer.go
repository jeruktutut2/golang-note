package setup

import (
	"golang-note/consumer"

	"github.com/rabbitmq/amqp091-go"
)

func RabbitmqConsummer(channel *amqp091.Channel) {
	consumer.ConsumeEmail(channel)
	consumer.ConsumerTextMessage(channel)
}

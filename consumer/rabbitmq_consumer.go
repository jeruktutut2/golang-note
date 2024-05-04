package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

func ConsumeEmail(channel *amqp091.Channel) {
	go func() {
		rabbitMqCtx := context.Background()
		emailConsumers, err := channel.ConsumeWithContext(rabbitMqCtx, "email", "email-consumer", true, false, false, false, nil)
		if err != nil {
			log.Fatalln("error when consuming email:", err)
		}
		println(time.Now().String(), "rabbitmq: email-consumer listening to email queue")
		for emailConsumer := range emailConsumers {
			// fmt.Println("emailConsumer.RoutingKey:", emailConsumer.RoutingKey)
			// fmt.Println("emailConsumer.Body:", string(emailConsumer.Body))
			fmt.Println("emailConsumer:", emailConsumer.RoutingKey, string(emailConsumer.Body))
		}
	}()
}

func ConsumerTextMessage(channel *amqp091.Channel) {
	go func() {
		ctx := context.Background()
		consumers, err := channel.ConsumeWithContext(ctx, "text-message", "text-message-consumer", true, false, false, false, nil)
		if err != nil {
			log.Fatalln("error when consuming email:", err)
		}
		println(time.Now().String(), "rabbitmq: text-message-consumer listening to text-message queue")
		for consumer := range consumers {
			fmt.Println("consumer:", consumer.RoutingKey, string(consumer.Body))
		}
	}()
}

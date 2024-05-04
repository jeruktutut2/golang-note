package main

import (
	"context"
	"golang-note/configuration"
	"golang-note/consumer"
	"golang-note/setup"
	"golang-note/util"
	"os"
	"os/signal"
)

func main() {
	config := configuration.NewConfiguration()

	mysql := util.NewMysqlConnection(config.MysqlUsername, config.MysqlPassword, config.MysqlHost, config.MysqlPort, config.MysqlDatabase, config.MysqlMaxOpenConnection, config.MysqlMaxIdleConnection, config.MysqlConnectionMaxLifetime, config.MysqlConnectionMaxIdletime)
	defer mysql.Close()

	redis := util.NewRedisConnection(config.RedisHost, config.RedisPort, config.RedisDatabase)
	defer redis.Close()

	e := setup.Echo(config)

	setup.Globals()

	validate := setup.Validate()

	// kafka
	emailEmitter := util.NewEmmiter([]string{"localhost:9092"}, "email")
	defer emailEmitter.Finish()
	textMessageEmitter := util.NewEmmiter([]string{"localhost:9092"}, "text-message")
	defer textMessageEmitter.Finish()

	// rabbit mq
	rabbitmqConnection, rabbitmqChannel := util.NewRabbitmqConneciton(config)
	defer rabbitmqConnection.Close()

	repository := setup.Repository()
	service := setup.Service(mysql, redis, validate, config, repository)
	controller := setup.Controller(emailEmitter, textMessageEmitter, rabbitmqChannel, service)
	setup.Route(e, config, redis, controller)

	setup.RabbitmqConsummer(rabbitmqChannel)

	ctxKafka, cancelKafka := context.WithCancel(context.Background())
	consumer.ConsumeEmailKafka(ctxKafka, []string{"localhost:9092"}, "email", "email-consumer-group")
	consumer.ConsumeTextMessageKafka(ctxKafka, []string{"localhost:9092"}, "text-message", "text-message-consumer-group")

	setup.StartEcho(e, config)
	defer setup.StopEcho(e)

	// wait
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
	cancelKafka()
}

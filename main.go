package main

import (
	"context"
	"golang-note/configuration"
	"golang-note/controller"
	"golang-note/exception"
	"golang-note/globals"
	"golang-note/helper"
	"golang-note/middleware"
	"golang-note/repository"
	"golang-note/route"
	"golang-note/service"
	"golang-note/util"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

func main() {
	config := configuration.NewConfiguration()

	mysql := util.NewMysqlConnection(config.MysqlUsername, config.MysqlPassword, config.MysqlHost, config.MysqlPort, config.MysqlDatabase, config.MysqlMaxOpenConnection, config.MysqlMaxIdleConnection, config.MysqlConnectionMaxLifetime, config.MysqlConnectionMaxIdletime)
	defer mysql.Close()

	redis := util.NewRedisConnection(config.RedisHost, config.RedisPort, config.RedisDatabase)
	defer redis.Close()

	e := echo.New()
	e.Use(echomiddleware.Recover())
	e.HTTPErrorHandler = exception.CustomHTTPErrorHandler
	e.Use(middleware.SetGlobalRequestId)
	e.Use(middleware.SetGlobalTimeout(config.ApplicationTimeout))
	e.Use(middleware.SetGlobalRequestLog)

	globals.Session = make(map[string]string)

	validate := validator.New()
	helper.UsernameValidator(validate)
	helper.PasswordValidator(validate)
	helper.TelephoneValidator(validate)

	// kafka
	emmiter := util.NewEmmiter()
	defer emmiter.Finish()

	userRepository := repository.NewUserRepository()
	permissionRepository := repository.NewPermissionRepository()
	userPermissionRepository := repository.NewUserPermissionRepository()
	redisRepository := repository.NewRedisRepository()
	contextRepository := repository.NewContextRepository()

	userService := service.NewUserService(mysql, redis, validate, config.JwtKey, config.JwtAccessTokenExpireTime, config.JwtRefreshTokenExpireTime, userRepository, permissionRepository, userPermissionRepository, redisRepository)
	contextService := service.NewContextService(mysql.GetDb(), contextRepository)
	pdfService := service.NewPdfService()
	streamJsonService := service.NewStreamJsonService()

	userController := controller.NewUserController(userService)
	restController := controller.NewRestController()
	contextController := controller.NewContextController(contextService)
	permissionController := controller.NewPermissionController()
	pdfController := controller.NewPdfController(pdfService)
	kafkaController := controller.NewKafkaController(emmiter)
	streamJsonController := controller.NewStreamJsonController(streamJsonService)
	serverSentEventController := controller.NewServerSentEventController()

	route.UserRoute(e, config.JwtKey, redis.GetClient(), userController)
	route.RestRoute(e, restController)
	route.ContextRoute(e, contextController)
	route.PermissionRoute(e, config.JwtKey, redis.GetClient(), permissionController)
	route.PdfRoute(e, pdfController)
	route.KafkaRoute(e, kafkaController)
	route.StreamJsonRoute(e, streamJsonController)
	route.ServerSentEventRoute(e, serverSentEventController)

	// kafka
	p, err := util.NewProcessor()
	if err != nil {
		log.Fatalln("error creating processor: ", err)
	}
	ctxKafka, cancelKafka := context.WithCancel(context.Background())
	doneKafka := make(chan bool)
	go func() {
		defer close(doneKafka)
		if err = p.Run(ctxKafka); err != nil {
			log.Fatalf("error running processor: %v", err)
		} else {
			log.Printf("Processor shutdown cleanly")
		}
	}()

	// echo server
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		if err := e.Start(":" + strconv.Itoa(int(config.ApplicationPort))); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// kafka
	waitKafka := make(chan os.Signal, 1)
	signal.Notify(waitKafka, syscall.SIGINT, syscall.SIGTERM)
	<-waitKafka
	cancelKafka()
	<-doneKafka

	// echo server
	<-ctx.Done()

	mysql.Close()
	println("mysql closed properly")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	println("server shutdown properly")
}

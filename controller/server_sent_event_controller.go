package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-note/globals"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type ServerSentEventController interface {
	HandleSSE(c echo.Context) error
	SendMessage(c echo.Context) error
	HandleSSEWithoutChannel(c echo.Context) error
	SendMessageWithoutChannel(c echo.Context) error
}

type ServerSentEventControllerImplementation struct {
}

func NewServerSentEventController() ServerSentEventController {
	return &ServerSentEventControllerImplementation{}
}

var messageChan chan string

func (controller *ServerSentEventControllerImplementation) HandleSSE(c echo.Context) error {
	// set timeout to 0
	c.SetRequest(c.Request().WithContext(context.Background()))
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	messageChan = make(chan string)

	defer func() {
		close(messageChan)
		messageChan = nil
		log.Printf("client connection close")
	}()

	flusher, _ := c.Response().Writer.(http.Flusher)

	for {

		select {
		case message := <-messageChan:
			// fmt.Println("message:", message)
			json.NewEncoder(c.Response()).Encode(message)
			// c.Response().Flush()
			flusher.Flush()
		case <-c.Request().Context().Done():
			fmt.Println("connection close")
			return nil
		}
	}
}

func (controller *ServerSentEventControllerImplementation) SendMessage(c echo.Context) error {
	message := c.Param("message")
	messageChan <- message
	return nil
}

func (controller *ServerSentEventControllerImplementation) HandleSSEWithoutChannel(c echo.Context) error {
	// c.SetRequest(c.Request().WithContext(context.Background()))
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")

	sseClient := globals.SSEClient{}
	sseClient.Id = time.Now().UnixMilli()
	sseClient.Context = c
	globals.SSEClients = append(globals.SSEClients, sseClient)

	<-c.Request().Context().Done()

sseClientsLoop:
	for i := 0; i < len(globals.SSEClients); i++ {
		if globals.SSEClients[i].Id == sseClient.Id {
			globals.SSEClients = append(globals.SSEClients[:i], globals.SSEClients[i+1:]...)
			break sseClientsLoop
		}
	}
	return nil
}

func (controller *ServerSentEventControllerImplementation) SendMessageWithoutChannel(c echo.Context) error {
	message := c.Param("message")
	for i := 0; i < len(globals.SSEClients); i++ {
		ctx := globals.SSEClients[i].Context
		json.NewEncoder(ctx.Response()).Encode(message)
		ctx.Response().Flush()
	}
	return c.JSON(http.StatusOK, "ok")
}

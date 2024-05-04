package controller

import (
	"encoding/json"
	"fmt"
	"golang-note/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StreamJsonController interface {
	Stream(c echo.Context) error
	StreamWithoutChannel(c echo.Context) error
	StreamWithSSE(c echo.Context) error
}

type StreamJsonControllerImplementation struct {
	StreamJsonService service.StreamJsonService
}

func NewStreamJsonController(streamJsonService service.StreamJsonService) StreamJsonController {
	return &StreamJsonControllerImplementation{
		StreamJsonService: streamJsonService,
	}
}

func (controller *StreamJsonControllerImplementation) Stream(c echo.Context) error {
	// cannot do it when using json, just text html
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)

	channelStream := make(chan string)
	controller.StreamJsonService.Stream(channelStream)
	enc := json.NewEncoder(c.Response())
	for value := range channelStream {
		fmt.Println("value:", value)
		if err := enc.Encode(value); err != nil {
			return err
		}
		c.Response().Flush()
		// time.Sleep(1 * time.Second)
	}
	return nil
}

func (controller *StreamJsonControllerImplementation) StreamWithoutChannel(c echo.Context) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)
	controller.StreamJsonService.StreamWithoutChannel(c)
	return nil
}

func (controller *StreamJsonControllerImplementation) StreamWithSSE(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	controller.StreamJsonService.StreamWithSSE(c)
	message := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	return c.JSON(http.StatusOK, message)
}

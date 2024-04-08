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

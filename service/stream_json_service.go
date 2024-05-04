package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

type StreamJsonService interface {
	Stream(streamChannel chan string) string
	StreamWithoutChannel(c echo.Context) string
	StreamWithSSE(c echo.Context) string
}

type StreamJsonServiceImplementation struct {
}

func NewStreamJsonService() StreamJsonService {
	return &StreamJsonServiceImplementation{}
}

func (service *StreamJsonServiceImplementation) Stream(streamChannel chan string) string {
	group := &sync.WaitGroup{}
	group.Add(1)
	// err := errors.New("error")
	go func() {
		defer group.Done()
		streamChannel <- "response 1"
		fmt.Println("streamChannel <- response 1")
		time.Sleep(2 * time.Second)

		streamChannel <- "response 2"
		fmt.Println("streamChannel <- response 2")
		time.Sleep(2 * time.Second)
		// err = errors.New("error")
		// if err != nil {
		// 	return
		// }
		streamChannel <- "response 3"
		fmt.Println("streamChannel <- response 3")
		time.Sleep(2 * time.Second)

		streamChannel <- "response 4"
		fmt.Println("streamChannel <- response 4")
		time.Sleep(2 * time.Second)

		streamChannel <- "response 5"
		fmt.Println("streamChannel <- response 5")
		time.Sleep(2 * time.Second)
	}()
	go func() {
		group.Wait()
		close(streamChannel)
	}()
	return "service"
}

func (service *StreamJsonServiceImplementation) StreamWithoutChannel(c echo.Context) string {
	json.NewEncoder(c.Response()).Encode("stream1")
	c.Response().Flush()
	time.Sleep(2 * time.Second)
	json.NewEncoder(c.Response()).Encode("stream2")
	c.Response().Flush()
	time.Sleep(2 * time.Second)
	json.NewEncoder(c.Response()).Encode("stream3")
	c.Response().Flush()
	return "stream"
}

func (service *StreamJsonServiceImplementation) StreamWithSSE(c echo.Context) string {
	json.NewEncoder(c.Response()).Encode("stream1")
	c.Response().Flush()
	time.Sleep(2 * time.Second)
	json.NewEncoder(c.Response()).Encode("stream2")
	c.Response().Flush()
	time.Sleep(2 * time.Second)
	json.NewEncoder(c.Response()).Encode("stream3")
	c.Response().Flush()
	return "stream"
}

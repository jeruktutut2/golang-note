package service

import (
	"fmt"
	"sync"
	"time"
)

type StreamJsonService interface {
	Stream(streamChannel chan string) string
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

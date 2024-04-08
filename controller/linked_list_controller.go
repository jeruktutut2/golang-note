package controller

import (
	"fmt"
	"net/http"

	"golang-note/middleware"
	"golang-note/service"
)

type LinkedListController interface {
	Singly(writer http.ResponseWriter, request *http.Request)
}

type LinkedListControllerImplementation struct {
	LinkedListService service.LinkedListService
}

func NewLinkedListController(linkedListService service.LinkedListService) LinkedListController {
	return &LinkedListControllerImplementation{
		LinkedListService: linkedListService,
	}
}

func (controller *LinkedListControllerImplementation) Singly(writer http.ResponseWriter, request *http.Request) {
	requestId := request.Context().Value(middleware.RequestIdKey).(string)
	controller.LinkedListService.Singly(request.Context(), requestId)
	fmt.Println("Singly")
}

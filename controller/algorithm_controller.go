package controller

import (
	"net/http"

	"golang-note/middleware"
	modelresponse "golang-note/model/response"
	"golang-note/service"

	"github.com/labstack/echo/v4"
)

type AlgorithmController interface {
	LinearSearch(c echo.Context) error
	BinarySearch(c echo.Context) error
	InterpolationSearch(c echo.Context) error
	JumpSearch(c echo.Context) error
	TernarySearch(c echo.Context) error
}

type AlgorithmControllerImplementation struct {
	AlgorithmService service.AlgorithmService
}

func NewAlgorithmController(algorithmService service.AlgorithmService) AlgorithmController {
	return &AlgorithmControllerImplementation{
		AlgorithmService: algorithmService,
	}
}

func (controller *AlgorithmControllerImplementation) LinearSearch(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arrvalue := controller.AlgorithmService.LinearSearch(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arrvalue, "")
}

func (controller *AlgorithmControllerImplementation) BinarySearch(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arrvalue := controller.AlgorithmService.BinarySearch(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arrvalue, "")
}

func (controller *AlgorithmControllerImplementation) InterpolationSearch(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arrvalue := controller.AlgorithmService.InterpolationSearch(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arrvalue, "")
}

func (controller *AlgorithmControllerImplementation) JumpSearch(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arrvalue := controller.AlgorithmService.JumpSearch(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arrvalue, "")
}

func (controller *AlgorithmControllerImplementation) TernarySearch(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arrvalue := controller.AlgorithmService.TernarySearch(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arrvalue, "")
}

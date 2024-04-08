package controller

import (
	"net/http"

	"golang-note/middleware"
	modelresponse "golang-note/model/response"
	"golang-note/service"

	"github.com/labstack/echo/v4"
	// "github.com/julienschmidt/httprouter"
)

type ArrayController interface {
	ReverseArray(c echo.Context) error
	RotationArray(c echo.Context) error
	RearrangeArray(c echo.Context) error
	RangeSumArray(c echo.Context) error
	RangeWithUpdateArray(c echo.Context) error
	// sparsetable pending
	MetricArray1(c echo.Context) error
	MetricArray2(c echo.Context) error
	MultiplyMatrix(c echo.Context) error
	KadanesAlgorithm(c echo.Context) error
	DutchNationalFlagAlgorithm(c echo.Context) error
}

type ArrayControllerImplementation struct {
	ArrayService service.ArrayService
}

func NewArrayController(arrayService service.ArrayService) ArrayController {
	return &ArrayControllerImplementation{
		ArrayService: arrayService,
	}
}

func (controller *ArrayControllerImplementation) ReverseArray(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.ArrayService.ReverseArray(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

func (controller *ArrayControllerImplementation) RotationArray(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.ArrayService.RotationArray(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

func (controller *ArrayControllerImplementation) RearrangeArray(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.ArrayService.RearrangeArray(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

func (controller *ArrayControllerImplementation) RangeSumArray(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.ArrayService.RangeSumArray(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

func (controller *ArrayControllerImplementation) RangeWithUpdateArray(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.ArrayService.RangeWithUpdateArray(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

func (controller *ArrayControllerImplementation) MetricArray1(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.ArrayService.MetricArray1(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

func (controller *ArrayControllerImplementation) MetricArray2(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.ArrayService.MetricArray2(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

func (controller *ArrayControllerImplementation) MultiplyMatrix(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.ArrayService.MultiplyMatrix(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

func (controller *ArrayControllerImplementation) KadanesAlgorithm(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	result := controller.ArrayService.KadanesAlgorithm(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, result, "")
}

func (controller *ArrayControllerImplementation) DutchNationalFlagAlgorithm(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.ArrayService.DutchNationalFlagAlgorithm(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

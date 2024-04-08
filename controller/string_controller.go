package controller

import (
	"net/http"

	"golang-note/middleware"
	modelresponse "golang-note/model/response"
	"golang-note/service"

	"github.com/labstack/echo/v4"
)

type StringController interface {
	Substring2(c echo.Context) error
	Substring1(c echo.Context) error
	Subsequence1(c echo.Context) error
	Subsequence2(c echo.Context) error
	Rotation(c echo.Context) error
	BinaryString(c echo.Context) error
	Palindrome(c echo.Context) error
	LexicographicRankString(c echo.Context) error
	PatternSearching(c echo.Context) error
}

type StringControllerImplementation struct {
	StringService service.StringService
}

func NewStringController(stringService service.StringService) StringController {
	return &StringControllerImplementation{
		StringService: stringService,
	}
}

func (controller *StringControllerImplementation) Substring1(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	numberOfContain := controller.StringService.Substring1(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, numberOfContain, "")
}

func (controller *StringControllerImplementation) Substring2(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.StringService.Substring2(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

func (controller *StringControllerImplementation) Subsequence1(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.StringService.Subsequence1(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

func (controller *StringControllerImplementation) Subsequence2(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	gks := controller.StringService.Subsequence2(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, gks, "")
}

func (controller *StringControllerImplementation) Rotation(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	srotation := controller.StringService.Rotation(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, srotation, "")
}

func (controller *StringControllerImplementation) BinaryString(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	s := controller.StringService.BinaryString(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, s, "")
}

func (controller *StringControllerImplementation) Palindrome(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	s := controller.StringService.Palindrome(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, s, "")
}

func (controller *StringControllerImplementation) LexicographicRankString(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	rank := controller.StringService.LexicographicRackString(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, rank, "")
}

func (controller *StringControllerImplementation) PatternSearching(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	arr := controller.StringService.PatternSearching(c.Request().Context(), requestId)
	return modelresponse.ToResponse(c, http.StatusOK, requestId, arr, "")
}

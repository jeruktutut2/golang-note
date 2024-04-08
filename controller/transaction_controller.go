package controller

import (
	"net/http"
	"time"

	"golang-note/exception"
	"golang-note/helper"
	"golang-note/middleware"
	modelrequest "golang-note/model/request"
	modelresponse "golang-note/model/response"
	"golang-note/service"

	"github.com/labstack/echo/v4"
)

type TransactionController interface {
	Purchase(c echo.Context) error
}

type TransactionControllerImplementation struct {
	TransactionService service.TransactionService
}

func NewTransactionController(transactionService service.TransactionService) TransactionController {
	return &TransactionControllerImplementation{
		TransactionService: transactionService,
	}
}

func (controller *TransactionControllerImplementation) Purchase(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	idKey := c.Request().Context().Value(middleware.IdKey).(float64)
	id := int32(idKey)
	purchaseTransactionRequest := modelrequest.PurchaseTransactionRequest{}
	if err := c.Bind(&purchaseTransactionRequest); err != nil {
		helper.PrintLogToTerminal(err, requestId)
		err = exception.NewBadRequestException(err.Error())
		return exception.ErrorHandler(c, requestId, err)
	}
	now := time.Now().Unix()
	purchaseResponse, err := controller.TransactionService.Purchase(c.Request().Context(), requestId, id, now, purchaseTransactionRequest)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}
	return modelresponse.ToResponse(c, http.StatusOK, requestId, purchaseResponse, "")
}

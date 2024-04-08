package controller

import (
	"fmt"
	"golang-note/exception"
	"golang-note/middleware"
	"golang-note/service"

	"github.com/labstack/echo/v4"
)

type PdfController interface {
	GetPdf(c echo.Context) error
}

type PdfControllerImplementation struct {
	PdfService service.PdfService
}

func NewPdfController(pdfService service.PdfService) PdfController {
	return &PdfControllerImplementation{
		PdfService: pdfService,
	}
}

func (controller *PdfControllerImplementation) GetPdf(c echo.Context) error {
	requestId := c.Request().Context().Value(middleware.RequestIdKey).(string)
	pdfBytes, err := controller.PdfService.GetPdf(c.Request().Context(), requestId)
	if err != nil {
		return exception.ErrorHandler(c, requestId, err)
	}
	c.Response().Header().Set("Content-Type", "application/pdf")
	c.Response().Header().Set("Content-Disposition", "attachment; filename=output.pdf")
	fmt.Println("pdfBytes.Bytes():", pdfBytes.Bytes())
	c.Response().Write(pdfBytes.Bytes())
	return nil
}

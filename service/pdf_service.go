package service

import (
	"bytes"
	"context"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
)

type PdfService interface {
	GetPdf(ctx context.Context, requestId string) (bytesBuffer bytes.Buffer, err error)
}

type PdfServiceImplementation struct {
}

func NewPdfService() PdfService {
	return &PdfServiceImplementation{}
}

// https://www.youtube.com/watch?v=jwOy4JgleTU
func (service *PdfServiceImplementation) GetPdf(ctx context.Context, requestId string) (bytesBuffer bytes.Buffer, err error) {
	maroto := pdf.NewMaroto(consts.Portrait, consts.A4)
	maroto.SetPageMargins(20, 10, 20)
	// only produces 1 output, if use two output, the other one will produce zero byte
	// err = maroto.OutputFileAndClose("pdfs/pdfA2.pdf")
	bytesBuffer, err = maroto.Output()
	return
}

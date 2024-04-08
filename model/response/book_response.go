package modelresponse

import (
	"errors"

	modelentity "golang-note/model/entity"
)

type GetBookByIdResponse struct {
	Id    int32
	Name  string
	Price float64
	Stock int16
}

func ToGetBookByIdResponse(book modelentity.Book) (getBookByIdResponse GetBookByIdResponse, err error) {
	getBookByIdResponse.Id = book.Id.Int32
	getBookByIdResponse.Name = book.Name.String

	var price float64
	price, ok := book.Price.Decimal.Float64()
	if !ok {
		getBookByIdResponse = GetBookByIdResponse{}
		err = errors.New("cannot convert from decimal to float64")
		return
	}
	getBookByIdResponse.Price = price
	getBookByIdResponse.Stock = book.Stock.Int16
	return
}

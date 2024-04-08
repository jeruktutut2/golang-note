package modelrequest

type BookPurchaseTransactionRequest struct {
	BookId   int32 `json:"bookId" validate:"required"`
	Quantity int16 `json:"quantity" validate:"required"`
}
type PurchaseTransactionRequest struct {
	Name                            string                           `json:"name" validate:"required"`
	BookPurchaseTransactionRequests []BookPurchaseTransactionRequest `json:"books" validate:"required,dive"`
}

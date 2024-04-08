package modelresponse

type PurchaseTransactionResponse struct {
	Id int32
}

func ToPurchaseTransactionResponse(id int32) (purchaseTransactionResponse PurchaseTransactionResponse) {
	// var purchaseTransactionResponseR PurchaseTransactionResponse
	// purchaseTransactionResponseR.Id = id
	// purchaseTransactionResponse = purchaseTransactionResponseR
	purchaseTransactionResponse.Id = id
	return
}

package modelresponse

import (
	"errors"

	modelentity "golang-note/model/entity"
)

type GetUserWalletResponse struct {
	Balancce float64 `json:"balance"`
}

func ToUserWallet(wallet modelentity.Wallet) (getUserWalletResponse GetUserWalletResponse, err error) {
	var balance float64
	balance, ok := wallet.Balance.Decimal.Float64()
	if !ok {
		err = errors.New("cannot convert from decimal to float64")
		return
	}
	getUserWalletResponse.Balancce = balance
	return
}

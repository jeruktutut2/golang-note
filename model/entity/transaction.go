package modelentity

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type Transaction struct {
	Id            sql.NullInt32
	UserId        sql.NullInt32
	Username      sql.NullString
	UserEmail     sql.NullString
	WalletId      sql.NullInt32
	WalletUserId  sql.NullInt32
	WalletBalance decimal.NullDecimal
	Paid          sql.NullInt16
	CreatedAt     sql.NullInt64
}

package modelentity

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type TransactionDetail struct {
	Id            sql.NullInt32
	TransactionId sql.NullInt32
	BookId        sql.NullInt32
	BookName      sql.NullString
	BookPrice     decimal.NullDecimal
	Quantity      sql.NullInt16
	CreatedAt     sql.NullInt64
}

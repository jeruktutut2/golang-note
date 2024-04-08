package modelentity

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type Book struct {
	Id    sql.NullInt32
	Name  sql.NullString
	Price decimal.NullDecimal
	Stock sql.NullInt16
}

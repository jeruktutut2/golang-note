package modelentity

import (
	"database/sql"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	Id      sql.NullInt32
	UserId  sql.NullInt32
	Balance decimal.NullDecimal
}

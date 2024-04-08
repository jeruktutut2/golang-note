package repository

import (
	"context"
	"database/sql"

	modelentity "golang-note/model/entity"
)

type TransactionDetailRepository interface {
	CreateMany(tx *sql.Tx, ctx context.Context, transactionDetails *[]*modelentity.TransactionDetail) (rowsAffected int64, err error)
}

type TransactionDetailRepositoryImplementation struct {
}

func NewTransactionDetailRepository() TransactionDetailRepository {
	return &TransactionDetailRepositoryImplementation{}
}

// why did *[]*modelentity.TransactionDetail, to prevent moved to heap: transactionDetail line 30
func (repository *TransactionDetailRepositoryImplementation) CreateMany(tx *sql.Tx, ctx context.Context, transactionDetails *[]*modelentity.TransactionDetail) (rowsAffected int64, err error) {
	query := `INSERT INTO transaction_detail (transaction_id, book_id, book_name, book_price, quantity, created_at) 
			VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	var result sql.Result
	for _, transactionDetail := range *transactionDetails {
		result, err = stmt.ExecContext(ctx, &transactionDetail.TransactionId, &transactionDetail.BookId, &transactionDetail.BookName, &transactionDetail.BookPrice, &transactionDetail.Quantity, &transactionDetail.CreatedAt)
		if err != nil {
			rowsAffected = 0
			return
		}

		var ra int64
		ra, err = result.RowsAffected()
		if err != nil {
			rowsAffected = 0
			return
		}
		rowsAffected += ra
	}
	return
}

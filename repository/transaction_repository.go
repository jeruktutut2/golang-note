package repository

import (
	"context"
	"database/sql"

	modelentity "golang-note/model/entity"
)

type TransactionRepository interface {
	Create(tx *sql.Tx, ctx context.Context, transaction *modelentity.Transaction) (rowsAffected int64, lastInsertId int64, err error)
}

type TransactionRepositoryImplementation struct {
}

func NewTransactionRepository() TransactionRepository {
	return &TransactionRepositoryImplementation{}
}

func (repository *TransactionRepositoryImplementation) Create(tx *sql.Tx, ctx context.Context, transaction *modelentity.Transaction) (rowsAffected int64, lastInsertId int64, err error) {
	query := `INSERT INTO transaction (user_id, username, user_email, wallet_id, wallet_user_id,
    								wallet_balance, paid, created_at) 
			VALUES (?, ?, ?, ?, ?,
				?, ?, ?)`
	result, err := tx.ExecContext(ctx, query, &transaction.UserId, &transaction.Username, &transaction.UserEmail, &transaction.WalletId, &transaction.WalletUserId,
		&transaction.WalletBalance, &transaction.Paid, &transaction.CreatedAt)
	if err != nil {
		return
	}
	rowsAffected, err = result.RowsAffected()
	if err != nil {
		rowsAffected = 0
		return
	}
	lastInsertId, err = result.LastInsertId()
	if err != nil {
		lastInsertId = 0
		return
	}
	return
}

package repository

import (
	"context"
	"database/sql"

	modelentity "golang-note/model/entity"
)

type WalletRepository interface {
	GetByUserId(db *sql.DB, ctx context.Context, wallet *modelentity.Wallet) (err error)
	GetByUserIdForUpdate(tx *sql.Tx, ctx context.Context, wallet *modelentity.Wallet) (err error)
	UpdateUserWalletBalance(tx *sql.Tx, ctx context.Context, wallet *modelentity.Wallet) (rowsAffected int64, err error)
}

type WalletRepositoryImplementation struct {
}

func NewWalletRepository() WalletRepository {
	return &WalletRepositoryImplementation{}
}

func (repository *WalletRepositoryImplementation) GetByUserId(db *sql.DB, ctx context.Context, wallet *modelentity.Wallet) (err error) {
	query := `SELECT id, user_id, balance FROM wallet WHERE user_id = ?;`
	err = db.QueryRowContext(ctx, query, &wallet.Id.Int32).Scan(&wallet.Id, &wallet.UserId, &wallet.Balance)
	if err != nil {
		wallet = nil
		return
	}
	return
}

func (repository *WalletRepositoryImplementation) GetByUserIdForUpdate(tx *sql.Tx, ctx context.Context, wallet *modelentity.Wallet) (err error) {
	query := `SELECT id, user_id, balance FROM wallet WHERE user_id = ? FOR UPDATE;`
	err = tx.QueryRowContext(ctx, query, &wallet.UserId.Int32).Scan(&wallet.Id, &wallet.UserId, &wallet.Balance)
	if err != nil {
		wallet = nil
		return
	}
	return
}

func (repository *WalletRepositoryImplementation) UpdateUserWalletBalance(tx *sql.Tx, ctx context.Context, wallet *modelentity.Wallet) (rowsAffected int64, err error) {
	query := `UPDATE wallet SET balance = ? WHERE id = ?;`
	result, err := tx.ExecContext(ctx, query, &wallet.Balance, &wallet.Id)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

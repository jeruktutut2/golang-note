package repository

import (
	"context"
	"database/sql"
	modelentity "golang-note/model/entity"
)

type UserRepository interface {
	GetByUsernameForUpdate(tx *sql.Tx, ctx context.Context, user *modelentity.User) (err error)
	GetByUsername(tx *sql.Tx, ctx context.Context, user *modelentity.User) (err error)
	GetByIdForUpdate(tx *sql.Tx, ctx context.Context, user *modelentity.User) (err error)
	FindByEmailForUpate(tx *sql.Tx, ctx context.Context, email string) (user modelentity.User, err error)
	UpdateRefreshToken(tx *sql.Tx, ctx context.Context, id int32, refreshToken string) (rowsAffected int64, err error)
	CountByRefreshToken(tx *sql.Tx, ctx context.Context, username string, refreshToken string, countRefreshToken *uint16) (err error)
	GetByRefreshToken(tx *sql.Tx, ctx context.Context, refreshToken string) (user modelentity.User, err error)
}

type UserRepositoryImplementaion struct {
}

func NewUserRepository() UserRepository {
	return &UserRepositoryImplementaion{}
}

func (repository *UserRepositoryImplementaion) GetByUsernameForUpdate(tx *sql.Tx, ctx context.Context, user *modelentity.User) (err error) {
	query := `SELECT id, username, email, password, refresh_token, utc, created_at FROM user WHERE username = ? FOR UPDATE;`
	err = tx.QueryRowContext(ctx, query, &user.Username.String).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.RefreshToken, &user.Utc, &user.CreatedAt)
	if err != nil {
		user = nil
		return
	}
	return
}

func (repository *UserRepositoryImplementaion) GetByUsername(tx *sql.Tx, ctx context.Context, user *modelentity.User) (err error) {
	query := `SELECT id, username, email, password, refresh_token, utc, created_at FROM user WHERE username = ?;`
	err = tx.QueryRowContext(ctx, query, &user.Username.String).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.RefreshToken, &user.Utc, &user.CreatedAt)
	if err != nil {
		user = nil
		return
	}
	return
}

func (repository *UserRepositoryImplementaion) GetByIdForUpdate(tx *sql.Tx, ctx context.Context, user *modelentity.User) (err error) {
	query := `SELECT id, username, email, password, refresh_token, utc, created_at FROM user WHERE id = ? FOR UPDATE;`
	err = tx.QueryRowContext(ctx, query, &user.Id.Int32).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.RefreshToken, &user.Utc, &user.CreatedAt)
	if err != nil {
		user = nil
		return
	}
	return
}

func (repository *UserRepositoryImplementaion) FindByEmailForUpate(tx *sql.Tx, ctx context.Context, email string) (user modelentity.User, err error) {
	query := `SELECT id, username, email, password, refresh_token, utc, created_at FROM user WHERE email = ? FOR UPDATE;`
	err = tx.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.RefreshToken, &user.Utc, &user.CreatedAt)
	if err != nil {
		return
	}
	return
}

func (repository *UserRepositoryImplementaion) UpdateRefreshToken(tx *sql.Tx, ctx context.Context, id int32, refreshToken string) (rowsAffected int64, err error) {
	query := `UPDATE user SET refresh_token = ? WHERE id = ?;`
	result, err := tx.ExecContext(ctx, query, refreshToken, id)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

func (repository *UserRepositoryImplementaion) CountByRefreshToken(tx *sql.Tx, ctx context.Context, username string, refreshToken string, countRefreshToken *uint16) (err error) {
	query := `SELECT COUNT(*) AS count_refresh_token FROM user WHERE username = ? AND refresh_token = ?;`
	err = tx.QueryRowContext(ctx, query, username, refreshToken).Scan(&countRefreshToken)
	if err != nil {
		countRefreshToken = nil
		return
	}
	return
}

func (repository *UserRepositoryImplementaion) GetByRefreshToken(tx *sql.Tx, ctx context.Context, refreshToken string) (user modelentity.User, err error) {
	query := `SELECT id, username, email, password, refresh_token, utc, created_at FROM user WHERE refresh_token = ?;`
	err = tx.QueryRowContext(ctx, query, refreshToken).Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.RefreshToken, &user.Utc, &user.CreatedAt)
	if err != nil {
		user = modelentity.User{}
		return
	}
	return
}

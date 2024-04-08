package repository

import (
	"context"
	"database/sql"

	modelentity "golang-note/model/entity"
)

type BookRepository interface {
	GetBookById(db *sql.DB, ctx context.Context, book *modelentity.Book) (err error)
	GetBookByIdForUpdate(tx *sql.Tx, ctx context.Context, book *modelentity.Book) (err error)
	UpdateStock(tx *sql.Tx, ctx context.Context, book *modelentity.Book) (rowsAffected int64, err error)
	GetBookByInIdsForUpdate(tx *sql.Tx, ctx context.Context, ids []int32, books *[]modelentity.Book) (err error)
	UpdateManyStock(tx *sql.Tx, ctx context.Context, books *[]modelentity.Book) (rowsAffected int64, err error)
}

type BookRepositoryImplementation struct {
}

func NewBookRepository() BookRepository {
	return &BookRepositoryImplementation{}
}

func (repository *BookRepositoryImplementation) GetBookById(db *sql.DB, ctx context.Context, book *modelentity.Book) (err error) {
	query := `SELECT id, name, price, stock FROM book WHERE id = ?;`
	err = db.QueryRowContext(ctx, query, &book.Id.Int32).Scan(&book.Id, &book.Name, &book.Price, &book.Stock)
	if err != nil {
		book = nil
		return
	}
	return
}

func (repository *BookRepositoryImplementation) GetBookByIdForUpdate(tx *sql.Tx, ctx context.Context, book *modelentity.Book) (err error) {
	query := `SELECT id, name, price, stock FROM book WHERE id = ? FOR UPDATE;`
	err = tx.QueryRowContext(ctx, query, &book.Id.Int32).Scan(&book.Id, &book.Name, &book.Price, &book.Stock)
	if err != nil {
		book = nil
		return
	}
	return
}

func (repository *BookRepositoryImplementation) UpdateStock(tx *sql.Tx, ctx context.Context, book *modelentity.Book) (rowsAffected int64, err error) {
	query := `UPDATE book SET stock = ? WHERE id = ?;`
	result, err := tx.ExecContext(ctx, query, &book.Stock, &book.Id)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

func (repository *BookRepositoryImplementation) GetBookByInIdsForUpdate(tx *sql.Tx, ctx context.Context, ids []int32, books *[]modelentity.Book) (err error) {
	var placeholder string
	var params []interface{}
	for _, id := range ids {
		placeholder += `,?`
		params = append(params, id)
	}
	placeholder = placeholder[1:]
	query := `SELECT id, name, price, stock FROM book WHERE id IN (` + placeholder + `) FOR UPDATE;`
	rows, err := tx.QueryContext(ctx, query, params...)
	if err != nil {
		return
	}

	defer func() {
		errRowsClose := rows.Close()
		if errRowsClose != nil {
			books = nil
			err = errRowsClose
		}
	}()

	for rows.Next() {
		var book modelentity.Book
		err = rows.Scan(&book.Id, &book.Name, &book.Price, &book.Stock)
		if err != nil {
			books = nil
			return
		}
		*books = append(*books, book)
	}
	return
}

func (repository *BookRepositoryImplementation) UpdateManyStock(tx *sql.Tx, ctx context.Context, books *[]modelentity.Book) (rowsAffected int64, err error) {
	query := `UPDATE book SET stock = ? WHERE id = ?;`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	var result sql.Result
	for _, book := range *books {
		result, err = stmt.ExecContext(ctx, &book.Stock, &book.Id)
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

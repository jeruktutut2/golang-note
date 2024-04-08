package repository

import (
	"context"
	"database/sql"
)

type ContextRepository interface {
	Timeout(db *sql.DB, ctx context.Context) (str string, err error)
	TimeoutTx(tx *sql.Tx, ctx context.Context) (str string, err error)
	CreateTable1(db *sql.DB, ctx context.Context) (rowsAffected int64, err error)
	CreateTable2(db *sql.DB, ctx context.Context) (rowsAffected int64, err error)
	CreateTable3(db *sql.DB, ctx context.Context) (rowsAffected int64, err error)
	CreateTable1Tx(tx *sql.Tx, ctx context.Context) (rowsAffected int64, err error)
	CreateTable2Tx(tx *sql.Tx, ctx context.Context) (rowsAffected int64, err error)
	CreateTable3Tx(tx *sql.Tx, ctx context.Context) (rowsAffected int64, err error)
}

type ContextRepositoryImplementation struct {
}

func NewContextRepository() ContextRepository {
	return &ContextRepositoryImplementation{}
}

func (repository *ContextRepositoryImplementation) Timeout(db *sql.DB, ctx context.Context) (str string, err error) {
	query := `SELECT SLEEP(60) AS sleep;`
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return
	}
	defer func() {
		errRowsClose := rows.Close()
		if errRowsClose != nil {
			str = ""
			err = errRowsClose
		}
	}()

	if rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			str = ""
			return
		}
	}
	return
}

func (repository *ContextRepositoryImplementation) TimeoutTx(tx *sql.Tx, ctx context.Context) (str string, err error) {
	query := `SELECT SLEEP(60) AS sleep;`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return
	}

	defer func() {
		errRowsClose := rows.Close()
		if errRowsClose != nil {
			str = ""
			err = errRowsClose
		}
	}()

	if rows.Next() {
		err = rows.Scan(&str)
		if err != nil {
			str = ""
			return
		}
	}
	return
}

func (repository *ContextRepositoryImplementation) CreateTable1(db *sql.DB, ctx context.Context) (rowsAffected int64, err error) {
	query := `INSERT INTO table1(table1) VALUES("table1");`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

func (repository *ContextRepositoryImplementation) CreateTable2(db *sql.DB, ctx context.Context) (rowsAffected int64, err error) {
	query := `INSERT INTO table2(table2) VALUES("table2");`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

func (repository *ContextRepositoryImplementation) CreateTable3(db *sql.DB, ctx context.Context) (rowsAffected int64, err error) {
	query := `INSERT INTO table3(table3) VALUES("table3");`
	result, err := db.ExecContext(ctx, query)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

func (repository *ContextRepositoryImplementation) CreateTable1Tx(tx *sql.Tx, ctx context.Context) (rowsAffected int64, err error) {
	query := `INSERT INTO table1(table1) VALUES("table1tx");`
	result, err := tx.ExecContext(ctx, query)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

func (repository *ContextRepositoryImplementation) CreateTable2Tx(tx *sql.Tx, ctx context.Context) (rowsAffected int64, err error) {
	query := `INSERT INTO table2(table2) VALUES("table2tx");`
	result, err := tx.ExecContext(ctx, query)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

func (repository *ContextRepositoryImplementation) CreateTable3Tx(tx *sql.Tx, ctx context.Context) (rowsAffected int64, err error) {
	query := `INSERT INTO table3(table3) VALUES("table3tx");`
	result, err := tx.ExecContext(ctx, query)
	if err != nil {
		return
	}
	return result.RowsAffected()
}

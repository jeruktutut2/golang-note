package repository

import (
	"context"
	"database/sql"
	modelentity "golang-note/model/entity"
)

type PermissionRepository interface {
	GetByInId(tx *sql.Tx, ctx context.Context, ids []interface{}, permissions *[]modelentity.Permission) (err error)
}

type PermissionRepositoryImplementation struct {
}

func NewPermissionRepository() PermissionRepository {
	return &PermissionRepositoryImplementation{}
}

func (repository *PermissionRepositoryImplementation) GetByInId(tx *sql.Tx, ctx context.Context, ids []interface{}, permissions *[]modelentity.Permission) (err error) {
	var placeholder string
	for i := 0; i < len(ids); i++ {
		placeholder += `,?`
	}
	placeholder = placeholder[1:]
	query := `SELECT id, permission FROM permission WHERE id IN (` + placeholder + `);`
	rows, err := tx.QueryContext(ctx, query, ids...)
	if err != nil {
		return
	}
	defer func() {
		errRowsClose := rows.Close()
		if errRowsClose != nil {
			permissions = nil
			err = errRowsClose
		}
	}()

	for rows.Next() {
		var permission modelentity.Permission
		err = rows.Scan(&permission.Id, &permission.Permission)
		if err != nil {
			permissions = nil
			return
		}
		*permissions = append(*permissions, permission)
	}
	return
}

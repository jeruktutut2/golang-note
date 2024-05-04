package repository

import (
	"context"
	"database/sql"
	modelentity "golang-note/model/entity"
)

type UserPermissionRepository interface {
	GetByUserId(tx *sql.Tx, ctx context.Context, userId int32) (userPermissions []modelentity.UserPermission, err error)
}

type UserPermissionRepositoryImplementation struct {
}

func NewUserPermissionRepository() UserPermissionRepository {
	return &UserPermissionRepositoryImplementation{}
}

func (repository *UserPermissionRepositoryImplementation) GetByUserId(tx *sql.Tx, ctx context.Context, userId int32) (userPermissions []modelentity.UserPermission, err error) {
	query := `SELECT id, user_id, permission_id FROM user_permission WHERE user_id = ?;`
	rows, err := tx.QueryContext(ctx, query, userId)
	if err != nil {
		return
	}
	defer func() {
		errRowsClose := rows.Close()
		if errRowsClose != nil {
			userPermissions = nil
			err = errRowsClose
		}
	}()

	for rows.Next() {
		var userPermission modelentity.UserPermission
		err = rows.Scan(&userPermission.Id, &userPermission.UserId, &userPermission.PermissionId)
		if err != nil {
			userPermissions = nil
			return
		}
		userPermissions = append(userPermissions, userPermission)
	}
	return
}

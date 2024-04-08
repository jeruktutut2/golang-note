package modelentity

import "database/sql"

type User struct {
	Id           sql.NullInt32
	Username     sql.NullString
	Email        sql.NullString
	Password     sql.NullString
	RefreshToken sql.NullString
	Utc          sql.NullString
	CreatedAt    sql.NullInt32
}

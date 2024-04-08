package modelentity

import "database/sql"

type Permission struct {
	Id         sql.NullInt16
	Permission sql.NullString
}

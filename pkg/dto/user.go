package dto

import (
	"database/sql"
	"encoding/json"
)

type DbBrokerUserData struct {
	UserId   sql.NullString
	Metadata json.RawMessage
}

type DBUserRef struct {
	UserId      sql.NullString `gorm:"column:user_id"`
	UserName    sql.NullString `gorm:"column:user_name"`
	Mobile      sql.NullString `gorm:"column:user_name"`
	Mpin        sql.NullString `gorm:"column:mpin"`
	CountryCode sql.NullString `gorm:"column:country_code"`
	UserRole    sql.NullString `gorm:"column:user_role"`
	CreatedAt   sql.NullTime   `gorm:"column:created_at"`
	UpdatedAt   sql.NullTime   `gorm:"column:updated_at"`
	DeletedAt   sql.NullTime   `gorm:"column:deleted_at"`
}

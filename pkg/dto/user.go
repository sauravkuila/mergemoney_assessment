package dto

import (
	"database/sql"
	"encoding/json"
)

type DbBrokerUserData struct {
	UserId   sql.NullString
	Metadata json.RawMessage
}

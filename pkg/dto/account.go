package dto

import (
	"database/sql"
	"time"
)

type AggregatorRequest struct {
	Mobile string `json:"mobile"`
}

type UserAccount struct {
	Type          string    `json:"type"`
	BankName      string    `json:"bank_name,omitempty"`
	AccountNumber string    `json:"account_number,omitempty"`
	Ifsc          string    `json:"ifsc,omitempty"`
	LinkedVia     string    `json:"linked_via"`
	WalletName    string    `json:"wallet_name,omitempty"`
	WalletID      string    `json:"wallet_id,omitempty"`
	UpiID         string    `json:"upi_id,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type GetAccountsResponse struct {
	Data []UserAccount `json:"data,omitempty"`
	CommonResponse
}

type DBUserAccount struct {
	Sid           sql.NullInt64  `gorm:"column:serialid;primaryKey"`
	UserId        sql.NullString `gorm:"column:userId;index"`
	Type          sql.NullString `gorm:"column:type"`
	BankName      sql.NullString `gorm:"column:bank_name"`
	AccountNumber sql.NullString `gorm:"column:account_number;index"`
	Ifsc          sql.NullString `gorm:"column:ifsc"`
	LinkedVia     sql.NullString `gorm:"column:linked_via"`
	WalletName    sql.NullString `gorm:"column:wallet_name"`
	WalletID      sql.NullString `gorm:"column:wallet_id;index"`
	UpiID         sql.NullString `gorm:"column:upi_id;index"`
	CreatedAt     sql.NullTime   `gorm:"column:created_at"`
	UpdatedAt     sql.NullTime   `gorm:"column:updated_at"`
	DeletedAt     sql.NullTime   `gorm:"column:deleted_at"`
}

func (obj *DBUserAccount) ToAggregatorAccount() UserAccount {
	return UserAccount{
		Type:          obj.Type.String,
		BankName:      obj.BankName.String,
		AccountNumber: obj.AccountNumber.String,
		Ifsc:          obj.Ifsc.String,
		LinkedVia:     obj.LinkedVia.String,
		WalletName:    obj.WalletName.String,
		WalletID:      obj.WalletID.String,
		UpiID:         obj.UpiID.String,
		CreatedAt:     obj.CreatedAt.Time,
		UpdatedAt:     obj.UpdatedAt.Time,
	}
}

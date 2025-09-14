package dto

import (
	"database/sql"
	"encoding/json"
)

type DBOrder struct {
	OrderID             sql.NullString  `gorm:"column:order_id;primaryKey"`
	UserID              sql.NullString  `gorm:"column:user_id;index"`
	SourceSID           sql.NullInt64   `gorm:"column:source_sid"`
	SourceCurrency      sql.NullString  `gorm:"column:source_currency"`
	SourceAmount        sql.NullFloat64 `gorm:"column:source_amount"`
	DestinationCurrency sql.NullString  `gorm:"column:destination_currency"`
	DestinationAmount   sql.NullFloat64 `gorm:"column:destination_amount"`
	ConversionRate      sql.NullFloat64 `gorm:"column:conversion_rate"`
	ConversionRateDate  sql.NullTime    `gorm:"column:conversion_rate_date"`
	OrderStatus         sql.NullString  `gorm:"column:order_status;index"`
	Remark              sql.NullString  `gorm:"column:remarks"`
	CreatedAt           sql.NullTime    `gorm:"column:created_at"`
	UpdatedAt           sql.NullTime    `gorm:"column:updated_at"`
}

type DBOrderDestination struct {
	DestinationID     sql.NullInt64  `gorm:"column:destination_id;primaryKey;autoIncrement"`
	OrderID           sql.NullString `gorm:"column:order_id;index"`
	DestinationType   sql.NullString `gorm:"column:destination_type;index"` // wallet | upi | netbanking
	WalletID          sql.NullString `gorm:"column:wallet_id"`              // for wallet
	UPIID             sql.NullString `gorm:"column:upi_id"`                 // for UPI
	BankAccountNumber sql.NullString `gorm:"column:bank_account_number"`    // for netbanking
	IFSCCode          sql.NullString `gorm:"column:ifsc_code"`              // for netbanking
	CreatedAt         sql.NullTime   `gorm:"column:created_at"`
}

type DBTransaction struct {
	TransactionID    sql.NullString  `gorm:"column:transaction_id;primaryKey"`
	OrderID          sql.NullString  `gorm:"column:order_id;index"`
	Provider         sql.NullString  `gorm:"column:provider;index"`    // e.g. Razorpay, Stripe
	ProviderID       sql.NullString  `gorm:"column:provider_id;index"` // provider’s reference number
	ProviderRequest  json.RawMessage `gorm:"column:provider_request;type:jsonb"`
	ProviderResponse json.RawMessage `gorm:"column:provider_response;type:jsonb"`
	Status           sql.NullString  `gorm:"column:status;index"`          // initiated → pending → inprogress → completed/failed
	ErrorMessage     sql.NullString  `gorm:"column:error_message"`         // error message if any
	RetryCount       int             `gorm:"column:retry_count;default:0"` // number of retries attempted
	LastRetryAt      sql.NullTime    `gorm:"column:last_retry_at"`         // timestamp of last retry
	CreatedAt        sql.NullTime    `gorm:"column:created_at"`
	UpdatedAt        sql.NullTime    `gorm:"column:updated_at"`
}

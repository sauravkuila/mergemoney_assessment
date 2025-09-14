package dto

type TransferRequest struct {
	Source      TransferRequestSource      `json:"source"`
	Destination TransferRequestDestination `json:"destination"`
}

type TransferRequestSource struct {
	Sid      int64   `json:"sid" binding:"required"`
	Currency string  `json:"currency" binding:"required,len=3"`
	Amount   float64 `json:"amount" binding:"required,min=0"`
}

type TransferRequestDestination struct {
	Currency        string                 `json:"currency" binding:"required,len=3"`
	RecipientDetail map[string]interface{} `json:"recipient_detail,omitempty"`
	Account         string                 `json:"account,omitempty"`
	SwiftCode       string                 `json:"swift_code,omitempty"`
	Upi             string                 `json:"upi,omitempty"`
	WalletID        string                 `json:"wallet_id,omitempty"`
}

type TransferResponse struct {
	Data *TransferData `json:"data,omitempty"`
	CommonResponse
}

// TODO
// make a factory to choose transfer vendor
// get rate of transfer
// save vendor and rate details in db
// return saved transaction id
// initiate payme nt once confirmed
// OR assume transfer rate is global, select vendor only to transfer

type TransferData struct {
	TransferID          string  `json:"transfer_id"`
	SourceCurrency      string  `json:"source_currency"`
	DestinationCurrency string  `json:"destination_currency"`
	SourceAmount        float64 `json:"source_amount"`
	DestinationAmount   float64 `json:"destination_amount"`
	ConversionRate      float64 `json:"conversion_rate"`
	ConversionRateDate  string  `json:"conversion_rate_date"`
	DestinationType     string  `json:"destination_type"`
}

type TransferConfirmRequest struct {
	TransferID string `json:"transfer_id" binding:"required"`
	Action     string `json:"action" binding:"required,oneof=confirm cancel"`
}

type TransferConfirmResponse struct {
	Data *TransferConfirm `json:"data,omitempty"`
	CommonResponse
}

type TransferConfirm struct {
	TransferID string `json:"transfer_id"`
	Status     string `json:"status"`
}

type TransferStatusRequest struct {
	TransferID string `uri:"transfer_id" binding:"required"`
}

type TransferStatusResponse struct {
	Data *TransferStatus `json:"data,omitempty"`
	CommonResponse
}

type TransferStatus struct {
	TransferID string `json:"transfer_id"`
	Status     string `json:"status"`
}

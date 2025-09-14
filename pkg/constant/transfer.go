package constant

const (
	TXN_INITIATED  = "initiated"
	TXN_PENDING    = "pending"
	TXN_INPROGRESS = "inprogress"
	TXN_COMPLETED  = "completed"
	TXN_FAILED     = "failed"
)

const (
	ORDER_CREATED    = "created"
	ORDER_INPROGRESS = "inprogress"
	ORDER_COMPLETED  = "completed"
	ORDER_FAILED     = "failed"
)

const (
	TRANSFER_DESTINATION_WALLET      = "wallet"
	TRANSFER_DESTINATION_UPI         = "upi"
	TRANSFER_DESTINATION_BANK        = "bank"
	TRANSFER_DESTINATION_CASH_PICKUP = "cash"
	TRANSFER_DESTINATION_UNKNOWN     = "unknown"
)

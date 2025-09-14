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
	TRANSFER_TYPE_WALLET      = "wallet"
	TRANSFER_TYPE_UPI         = "upi"
	TRANSFER_TYPE_BANK        = "bank"
	TRANSFER_TYPE_CASH_PICKUP = "cash"
	TRANSFER_TYPE_UNKNOWN     = "unknown"
)

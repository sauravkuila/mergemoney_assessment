package paymentprovider

type PayType string

const (
	PayTypeNetBanking PayType = "bank"
	PayTypeWallet     PayType = "wallet"
	PayTypeUPI        PayType = "upi"
	PayTypeCash       PayType = "cash"
)

type PaymentDetails struct {
	SourceDetail      PaymentInfo
	DestinationDetail PaymentInfo
	Remark            string
}

type PaymentInfo struct {
	Type          PayType
	AccountNumber string
	SwiftCode     string
	WalletID      string
	UPIID         string
	Cash          string
	Name          string
}

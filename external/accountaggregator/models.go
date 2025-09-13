package accountaggregator

type VendorResponse struct {
	Status      bool        `json:"status"`
	Description string      `json:"description"`
	Data        *VendorData `json:"data"`
}

type VendorData struct {
	Mobile   string `json:"mobile"`
	Accounts []struct {
		Type          string `json:"type"`
		BankName      string `json:"bank_name,omitempty"`
		AccountNumber string `json:"account_number,omitempty"`
		Ifsc          string `json:"ifsc,omitempty"`
		LinkedVia     string `json:"linked_via"`
		WalletName    string `json:"wallet_name,omitempty"`
		WalletID      string `json:"wallet_id,omitempty"`
		UpiID         string `json:"upi_id,omitempty"`
	} `json:"accounts"`
}

type VendorRequest struct {
	Mobile string `json:"mobile" binding:"required"`
}

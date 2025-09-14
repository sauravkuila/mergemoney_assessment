package accountaggregator

import (
	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"github.com/sauravkuila/mergemoney_assessment/pkg/utils"
)

func GetAccountsAgainstMobile(c *gin.Context, mobile string, util utils.UtilsItf) ([]dto.UserAccount, error) {
	var (
		accountData []dto.UserAccount
	)

	// TO DO: Integrate with actual account aggregator service
	// For now, returning a mock response
	// You can replace this with actual API calls to the account aggregator service

	// Mock response
	accountData = []dto.UserAccount{
		{
			Type:          "bank",
			BankName:      "HDFC Bank",
			AccountNumber: "XXXXXX1234",
			Ifsc:          "HDFC0001234",
			LinkedVia:     "netbanking",
		},
		{
			Type:       "wallet",
			WalletName: "Paytm",
			WalletID:   "paytm-9999999999",
			LinkedVia:  "mobile_number",
		},
		{
			Type:      "upi",
			UpiID:     "9999999999@okicici",
			LinkedVia: "upi",
		},
	}

	// Example of how the response would look like
	return accountData, nil
}

package paymentprovider

import (
	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"go.uber.org/zap"
)

var (
	PROVIDER_1_BASE_URL   string
	PROVIDER_1_SECRET_KEY string
	PROVIDER_2_BASE_URL   string
	PROVIDER_2_API_KEY    string
	providers             []string = []string{"provider_1", "provider_2"}
	providerConfigMap     map[string]map[string]string
)

func InitPaymentProvider() {
	// Initialize the Payment Provider package
	// This can include setting up any necessary configurations or connections
	PROVIDER_1_BASE_URL = config.GetConfig().GetString("external.paymentprovider.provider_1.url")
	PROVIDER_1_SECRET_KEY = config.GetConfig().GetString("external.paymentprovider.provider_1.secret_key")
	PROVIDER_2_BASE_URL = config.GetConfig().GetString("external.paymentprovider.provider_2.url")
	PROVIDER_2_API_KEY = config.GetConfig().GetString("external.paymentprovider.provider_2.api_key")

	providers = []string{"provider_1", "provider_2"}
	providerConfigMap = make(map[string]map[string]string)
	providerConfigMap["provider_1"] = map[string]string{
		"url":        PROVIDER_1_BASE_URL,
		"secret_key": PROVIDER_1_SECRET_KEY,
	}
	providerConfigMap["provider_2"] = map[string]string{
		"url":     PROVIDER_2_BASE_URL,
		"api_key": PROVIDER_2_API_KEY,
	}

	logger.Log().Info("Payment Provider initialized", zap.String("provider_1_url", PROVIDER_1_BASE_URL), zap.String("provider_2_url", PROVIDER_2_BASE_URL))
}

// type PaymentProviderItf interface {
// 	// Method to initiate a transfer
// 	InitiateTransfer(amount float64, currency string, payDetail PaymentDetails) (string, error)
// }

package paymentprovider

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"github.com/sauravkuila/mergemoney_assessment/pkg/utils"
	"go.uber.org/zap"
)

// returns transfer ID, provider used and error if any
func InitiateTransfer(ctx context.Context, amount float64, currency string, payDetail PaymentDetails, utilObj utils.UtilsItf) (string, string, json.RawMessage, json.RawMessage, error) {
	// Logic to choose a provider and initiate the transfer
	// fetch a provider preference from config
	provider := config.GetConfig().GetString("external.paymentprovider.preferred_provider")
	if provider == "" {
		provider = providers[0] // default to provider_1 if not set
	}
	providerConfig, exists := providerConfigMap[provider]
	if !exists {
		return "", "", nil, nil, errors.New("invalid payment provider") // or return an error indicating invalid provider
	}
	logger.Log().Info("Initiating transfer", zap.String("provider", provider), zap.Float64("amount", amount), zap.String("currency", currency), zap.Any("payment_details", payDetail), zap.Any("provider_config", providerConfig))

	// Here, you would add the logic to make an HTTP request to the selected provider's API
	// using the providerConfig for authentication and endpoint details.
	// For simplicity, we will skip the actual HTTP request and simulate a response.

	// Simulate a successful transfer initiation
	// In a real implementation, you would parse the response from the provider's API
	// and extract the transfer ID or handle errors accordingly.

	// For simplicity, we will just return a dummy transfer ID
	providerUniqueId := utilObj.GetUniqueId(provider) // this is provider specific idempotency key
	sampleproviderRequest, _ := json.Marshal(map[string]interface{}{
		"amount":        amount,
		"currency":      currency,
		"paymentDetail": payDetail,
	})
	sampleproviderResponse, _ := json.Marshal(map[string]interface{}{
		"status":      "success",
		"transfer_id": providerUniqueId,
	})
	logger.Log().Info("Transfer initiated with provider", zap.String("provider", provider), zap.String("provider_transfer_id", providerUniqueId), zap.Any("provider_request", string(sampleproviderRequest)), zap.Any("provider_response", string(sampleproviderResponse)))
	return providerUniqueId, provider, sampleproviderRequest, sampleproviderResponse, nil
}

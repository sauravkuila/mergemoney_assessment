package paymentprovider

import (
	"context"
	"errors"

	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"go.uber.org/zap"
)

// returns transfer ID, provider used and error if any
func InitiateTransfer(ctx context.Context, amount float64, currency string, payDetail PaymentDetails) (string, string, error) {
	// Logic to choose a provider and initiate the transfer
	// fetch a provider preference from config
	provider := config.GetConfig().GetString("external.paymentprovider.preferred_provider")
	if provider == "" {
		provider = providers[0] // default to provider_1 if not set
	}
	providerConfig, exists := providerConfigMap[provider]
	if !exists {
		return "", "", errors.New("invalid payment provider") // or return an error indicating invalid provider
	}
	logger.Log().Info("Initiating transfer", zap.String("provider", provider), zap.Float64("amount", amount), zap.String("currency", currency), zap.Any("payment_details", payDetail), zap.Any("provider_config", providerConfig))

	// Here, you would add the logic to make an HTTP request to the selected provider's API
	// using the providerConfig for authentication and endpoint details.
	// For simplicity, we will skip the actual HTTP request and simulate a response.

	// Simulate a successful transfer initiation
	// In a real implementation, you would parse the response from the provider's API
	// and extract the transfer ID or handle errors accordingly.

	// For simplicity, we will just return a dummy transfer ID
	return "dummy-transfer-id", provider, nil
}

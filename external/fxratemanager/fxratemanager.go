package fxratemanager

import (
	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"go.uber.org/zap"
)

var (
	FX_RATE_VENDOR_BASE_URL string
)

func InitFxRateVendor() {
	// Initialize the FX Rate Vendor package
	// This can include setting up any necessary configurations or connections
	FX_RATE_VENDOR_BASE_URL = config.GetConfig().GetString("external.fx_rate_vendor.url")
	logger.Log().Info("FX Rate Vendor initialized", zap.String("fx_rate_vendor_url", FX_RATE_VENDOR_BASE_URL))
}

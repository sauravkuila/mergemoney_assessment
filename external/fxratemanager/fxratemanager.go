package fxratemanager

import "github.com/sauravkuila/mergemoney_assessment/pkg/config"

var (
	FX_RATE_VENDOR_BASE_URL string
	FX_RATE_VENDOR_KEY      string
)

func InitFxRateVendor() {
	// Initialize the FX Rate Vendor package
	// This can include setting up any necessary configurations or connections
	FX_RATE_VENDOR_BASE_URL = config.GetConfig().GetString("external.fx_rate_vendor.url")
	FX_RATE_VENDOR_KEY = config.GetConfig().GetString("external.fx_rate_vendor.key")
}

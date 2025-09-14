package external

import (
	"github.com/sauravkuila/mergemoney_assessment/external/accountaggregator"
	"github.com/sauravkuila/mergemoney_assessment/external/fxratemanager"
	"github.com/sauravkuila/mergemoney_assessment/external/paymentprovider"
)

func InitExternal() {
	// Initialize the external package
	// This can include setting up any necessary configurations or connections
	accountaggregator.InitAccountAggregator()
	fxratemanager.InitFxRateVendor()
	paymentprovider.InitPaymentProvider()
}

//TODO: build a factory and an interface to dynamically match for providers

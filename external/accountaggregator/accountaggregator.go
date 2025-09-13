package accountaggregator

import "github.com/sauravkuila/mergemoney_assessment/pkg/config"

var (
	ACCOUNT_AGGREGATOR_BASE_URL string
	ACCOUNT_AGGREGATOR_KEY      string
)

func InitAccountAggregator() {
	// Initialize the AA package
	// This can include setting up any necessary configurations or connections
	ACCOUNT_AGGREGATOR_BASE_URL = config.GetConfig().GetString("external.account_aggregator.url")
	ACCOUNT_AGGREGATOR_KEY = config.GetConfig().GetString("external.account_aggregator.key")
}

package external

import (
	"github.com/sauravkuila/mergemoney_assessment/external/accountaggregator"
)

func InitExternal() {
	// Initialize the external package
	// This can include setting up any necessary configurations or connections
	accountaggregator.InitAccountAggregator()
}

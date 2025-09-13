package utils

import (
	"github.com/sauravkuila/mergemoney_assessment/pkg/utils/identifier"
	"github.com/sauravkuila/mergemoney_assessment/pkg/utils/restclient"
)

type UtilsItf interface {
	//returns a unique id everytime with appended argument
	identifier.IdentifierItf
	restclient.RestClientItf
}

type utilsSt struct {
	identifier.IdentifierItf
	restclient.RestClientItf
}

func GetUtilsObj(config UtilsConfig) UtilsItf {
	utilsSt := utilsSt{}
	utilsSt.IdentifierItf = identifier.GetIdentifierItf(config.Identifier)
	utilsSt.RestClientItf = restclient.GetRestClientItf(config.RestSSlEnabled)
	return utilsSt
}

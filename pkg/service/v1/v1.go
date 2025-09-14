package v1

import (
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao"
	"github.com/sauravkuila/mergemoney_assessment/pkg/service/v1/account"
	"github.com/sauravkuila/mergemoney_assessment/pkg/service/v1/login"
	"github.com/sauravkuila/mergemoney_assessment/pkg/service/v1/webhook"
)

type V1Group interface {
	login.LoginItf
	account.AccountItf
	webhook.ProviderWebhookItf
}

type v1 struct {
	login.LoginItf
	account.AccountItf
	webhook.ProviderWebhookItf
}

func GetV1Group(dao dao.RepositoryItf) V1Group {
	return &v1{
		login.GetLoginItf(dao),
		account.GetAccountItf(dao),
		webhook.GetProviderWebhookItf(dao),
	}
}

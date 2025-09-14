package webhook

import (
	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao"
)

type ProviderWebhookItf interface {
	ReconcileTransferFromProvider1(c *gin.Context)
	ReconcileTransferFromProvider2(c *gin.Context)
}

type providerWebhookSt struct {
	DB dao.RepositoryItf
}

func GetProviderWebhookItf(dao dao.RepositoryItf) ProviderWebhookItf {
	return &providerWebhookSt{
		DB: dao,
	}
}

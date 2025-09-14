package webhook

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (obj *providerWebhookSt) ReconcileTransferFromProvider2(c *gin.Context) {
	log.Println("webhook hit for provider 2")
}

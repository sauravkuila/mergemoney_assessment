package account

import (
	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao"
	"github.com/sauravkuila/mergemoney_assessment/pkg/utils"
	"github.com/sauravkuila/mergemoney_assessment/pkg/utils/identifier"
)

type AccountItf interface {
	// fetch user accounts
	GetAccounts(c *gin.Context)
	// initiate a transfer
	Transfer(c *gin.Context)
	// confirm or cancel a pending transfer
	TransferConfirm(c *gin.Context)
}

type accountSt struct {
	DB       dao.RepositoryItf
	utilsObj utils.UtilsItf
}

func GetAccountItf(dao dao.RepositoryItf) AccountItf {
	utilsConfig := utils.UtilsConfig{
		Identifier: identifier.IDENTIFIER_SNOWFLAKE,
	}
	return &accountSt{
		DB:       dao,
		utilsObj: utils.GetUtilsObj(utilsConfig),
	}
}

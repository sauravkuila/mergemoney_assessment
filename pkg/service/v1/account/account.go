package account

import (
	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao"
	"github.com/sauravkuila/mergemoney_assessment/pkg/utils"
)

type AccountItf interface {
	GetAccounts(c *gin.Context)
}

type accountSt struct {
	DB       dao.RepositoryItf
	utilsObj utils.UtilsItf
}

func GetAccountItf(dao dao.RepositoryItf) AccountItf {
	utilsConfig := utils.UtilsConfig{}
	return &accountSt{
		DB:       dao,
		utilsObj: utils.GetUtilsObj(utilsConfig),
	}
}

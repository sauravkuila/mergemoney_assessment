package login

import (
	"github.com/gin-gonic/gin"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao"
)

type LoginItf interface {
	GenerateOTP(c *gin.Context)
	VerifyOTP(c *gin.Context)
	SetMPIN(c *gin.Context)
	ResetMPIN(c *gin.Context)
	Refresh1FA(c *gin.Context)
	Refresh2FA(c *gin.Context)
	VerifyMPIN(c *gin.Context)
}

type loginSt struct {
	DB dao.RepositoryItf
	// utils beutils.UtilsItf
}

func GetLoginItf(dao dao.RepositoryItf) LoginItf {
	return &loginSt{
		DB: dao,
	}
}

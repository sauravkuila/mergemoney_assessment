package dao

import (
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao/account"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao/user"
	"gorm.io/gorm"
)

type RepositoryItf interface {
	user.DbUserItf
	account.DbAccountItf
}

type repositorySt struct {
	user.DbUserItf
	account.DbAccountItf
}

func GetRepositoryItf(psql *gorm.DB) RepositoryItf {
	return &repositorySt{
		user.GetUserItf(psql),
		account.GetAccountItf(psql),
	}
}

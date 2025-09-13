package dao

import (
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao/user"
	"gorm.io/gorm"
)

type RepositoryItf interface {
	user.DbUserItf
}

type repositorySt struct {
	user.DbUserItf
}

func GetRepositoryItf(psql *gorm.DB) RepositoryItf {
	return &repositorySt{
		user.GetUserItf(psql),
	}
}

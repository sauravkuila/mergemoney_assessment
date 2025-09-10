package dao

import (
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao/user"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type RepositoryItf interface {
	user.DbUserItf
}

type repositorySt struct {
	user.DbUserItf
}

func GetRepositoryItf(psql *gorm.DB, mnsql *mongo.Client) RepositoryItf {
	return &repositorySt{
		// user.GetUserItf(psql, mnsql),
	}
}

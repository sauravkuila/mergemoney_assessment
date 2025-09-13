package user

import (
	"context"

	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type DbUserItf interface {
	//fetch user data from number
	GetUserFromMobile(ctx context.Context, mobile string, countryCode string) (*dto.DBUserRef, error)
}

type userSt struct {
	psql  *gorm.DB
	mnsql *mongo.Client
}

func GetUserItf(psql *gorm.DB, mnsql *mongo.Client) DbUserItf {
	return &userSt{
		psql:  psql,
		mnsql: mnsql,
	}
}

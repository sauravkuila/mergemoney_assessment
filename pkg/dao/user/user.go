package user

import (
	"context"

	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"gorm.io/gorm"
)

type DbUserItf interface {
	//fetch user data from number
	GetUserFromMobile(ctx context.Context, mobile string, countryCode string) (*dto.DBUserRef, error)

	//fetch user data from user id
	GetUserFromUserId(ctx context.Context, userId string) (*dto.DBUserRef, error)

	//update user mpin
	SetUserMPIN(ctx context.Context, userId string, mpin string) error

	//reset any mpin
	ResetUserMPIN(ctx context.Context, mobile string, countryCode string) error
}

type userSt struct {
	psql *gorm.DB
}

func GetUserItf(psql *gorm.DB) DbUserItf {
	return &userSt{
		psql: psql,
	}
}

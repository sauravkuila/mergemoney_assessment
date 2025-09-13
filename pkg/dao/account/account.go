package account

import (
	"context"

	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"gorm.io/gorm"
)

type DbAccountItf interface {
	//fetch user data from number
	GetUserAccountsByUserId(ctx context.Context, userId string) ([]dto.DBUserAccount, error)

	// save user account
	SaveUserAccounts(ctx context.Context, userId string, account []dto.UserAccount) ([]dto.DBUserAccount, error)
}

type accountSt struct {
	psql *gorm.DB
}

func GetAccountItf(psql *gorm.DB) DbAccountItf {
	return &accountSt{
		psql: psql,
	}
}

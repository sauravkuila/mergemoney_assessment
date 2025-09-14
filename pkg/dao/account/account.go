package account

import (
	"context"

	"github.com/sauravkuila/mergemoney_assessment/pkg/dto"
	"gorm.io/gorm"
)

type DbAccountItf interface {
	//fetch user data from number
	GetUserAccountsByUserId(ctx context.Context, userId string) ([]dto.DBUserAccount, error)

	//fetch user account data from sid
	GetUserAccountsBySid(ctx context.Context, userId string, sid int64) (*dto.DBUserAccount, error)

	// save user account
	SaveUserAccounts(ctx context.Context, userId string, account []dto.UserAccount) ([]dto.DBUserAccount, error)

	// Save the order
	SaveOrder(ctx context.Context, order dto.DBOrder, orderDestination dto.DBOrderDestination) error

	// Update the order status
	UpdateOrderStatus(ctx context.Context, orderId string, status string, remark string) error

	// Get the order by id
	GetOrderById(ctx context.Context, orderId string, userId string) (*dto.DBOrder, *dto.DBOrderDestination, error)
}

type accountSt struct {
	psql *gorm.DB
}

func GetAccountItf(psql *gorm.DB) DbAccountItf {
	return &accountSt{
		psql: psql,
	}
}

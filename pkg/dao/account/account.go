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

	// Save the transaction
	SaveTransaction(ctx context.Context, transaction dto.DBTransaction) error

	// Update the transaction status
	UpdateTransactionStatus(ctx context.Context, providerId string, transactionId string, status string, remark string) error

	// Get the transaction by id
	GetTransactionById(ctx context.Context, transactionId string) (*dto.DBTransaction, error)

	// Get the transaction by order id
	GetTransactionByOrderId(ctx context.Context, orderId string) ([]dto.DBTransaction, error)
}

type accountSt struct {
	psql *gorm.DB
}

func GetAccountItf(psql *gorm.DB) DbAccountItf {
	return &accountSt{
		psql: psql,
	}
}

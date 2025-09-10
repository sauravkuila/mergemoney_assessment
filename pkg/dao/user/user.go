package user

import (
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type DbUserItf interface {
	// send data to be updated in metadata tag.
	// data is updated based on userId+brokerId combination OR based on brokerUserIdentifier.
	// preference of brokerUserIdentifier is taken if all of them are received.
	// UpdateUsersData(ctx context.Context, data dto.DbBrokerUserData) error

	// pass either userId+brokerId OR brokerUserIdentifier.
	// preference of brokerUserIdentifier is taken if all of them are received.
	// GetUserData(ctx context.Context, data dto.DbBrokerUserData) (*dto.DbBrokerUserData, error)

	//fetch user data from userId
	// GetUserFromUserId(ctx context.Context, userId string) (*dto.DBUserRef, error)
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

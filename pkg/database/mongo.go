package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func ConnectMongo(config DbConfig) (*mongo.Client, error) {
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s",
		config.User,
		config.Password,
		config.Host,
		config.Database,
	)
	if config.Port != 0 {
		mongoURI = fmt.Sprintf("mongodb://%s:%s@%s:%d/%s",
			config.User,
			config.Password,
			config.Host,
			config.Port,
			config.Database,
		)
	}

	clientOptions := options.Client().ApplyURI(mongoURI).
		SetRetryWrites(true).
		SetWriteConcern(writeconcern.Majority())

	//add default timeout if not provided
	timeout := config.ConnectTimeout
	if timeout <= 0 {
		timeout = DEFAULT_CONNECT_TIMEOUT
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	//connect to dsn
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Printf("failed to connect to mongoDB. err: %v. dsn: %s\n", err, mongoURI)
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	if config.ValidatePing {
		err = client.Ping(ctx, nil)
		if err != nil {
			defer client.Disconnect(ctx)
			fmt.Printf("failed to ping mongoDB. err: %v. dsn: %s\n", err, mongoURI)
			return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
		}
	}

	log.Println("mongo database connected")
	return client, nil
}

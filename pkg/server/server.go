package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/dao"
	"github.com/sauravkuila/mergemoney_assessment/pkg/database"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"github.com/sauravkuila/mergemoney_assessment/pkg/service"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	srv       *http.Server
	ctx       context.Context
	databases []*gorm.DB
)

// starts the server with initializations
func Start() error {
	ctx = context.Background()

	databases = make([]*gorm.DB, 0)

	logger.Log().Debug("config data",
		zap.Any("psql host", config.GetConfig().Get("databases.postgres.host")),
		zap.Any("psql port", config.GetConfig().Get("databases.postgres.port")),
		zap.Any("psql user", config.GetConfig().Get("databases.postgres.user")),
		zap.Any("psql password", config.GetConfig().Get("databases.postgres.password")),
		zap.Any("psql db", config.GetConfig().Get("databases.postgres.db")),
		zap.Any("psql sslmode", config.GetConfig().Get("databases.postgres.sslmode")),
		zap.Any("psql connect_timeout", config.GetConfig().Get("databases.postgres.connect_timeout")),
	)
	//connect to psql db
	psqlCfg := database.DbConfig{
		Host:           config.GetConfig().GetString("databases.postgres.host"),
		Port:           config.GetConfig().GetInt("databases.postgres.port"),
		User:           config.GetConfig().GetString("databases.postgres.user"),
		Password:       config.GetConfig().GetString("databases.postgres.password"),
		Database:       config.GetConfig().GetString("databases.postgres.db"),
		SSLMode:        config.GetConfig().GetString("databases.postgres.sslmode"),
		ConnectTimeout: config.GetConfig().GetInt("databases.postgres.connect_timeout"),
		PostgresConfig: database.PostgresConfig{
			Logger:           logger.GetDBLogger(),
			ApplicationName:  fmt.Sprintf("service-%s-%s-%d", os.Getenv("POD_NAME"), os.Getenv("POD_NAMESPACE"), os.Getpid()),
			StatementTimeout: 5,  //TODO: make this configurable
			MaxIdleConns:     10, //TODO: make this configurable
			MaxOpenConns:     30, //TODO: make this configurable
		},
	}
	postgresConn, err := database.ConnectPostgres(psqlCfg)
	if err != nil {
		logger.Log().Error("Failed to connect psql database", zap.Error(err))
		return err
	}
	databases = append(databases, postgresConn)

	// //init repo, factory, service and controller interfaces
	repo := dao.GetRepositoryItf(postgresConn)
	service := service.GetServiceItf(repo)

	// // create and start the gin server
	startRouter(service)

	return nil
}

// stops the router running in the go routine.
func ShutdownRouter() {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	logger.Log().Info("Shutting down router START")
	defer logger.Log().Info("Shutting down router END")
	if err := srv.Shutdown(timeoutCtx); err != nil {
		logger.Log().Fatal("Server forced to shutdown", zap.Error(err))
	}
	// catching ctx.Done(). timeout of 2 seconds.
	<-timeoutCtx.Done()
	log.Println("timeout of 2 seconds.")
}

// closes all database connections
//
//	closes each database connection made which are saved globally
//	logs error if unable to close
//		function used: *sql.DB.Close()
func CloseDatabase() {
	logger.Log().Info("disconnecting databases START")
	defer logger.Log().Info("disconnecting databases END")
	for _, database := range databases {
		db, _ := database.DB()
		if db != nil {
			err := db.Close()
			if err != nil {
				logger.Log().Error("unable to close db", zap.Error(err))
			}
		} else {
			logger.Log().Error("unable to close db as connection is nil")
		}
	}
}

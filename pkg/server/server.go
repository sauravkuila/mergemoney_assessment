package server

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var (
	srv            *http.Server
	ctx            context.Context
	databases      []*gorm.DB
	mongoDatabases []*mongo.Client
)

// starts the server with initializations
func Start() error {
	// ctx = context.Background()

	// databases = make([]*gorm.DB, 0)

	// belogger.Log().Debug("config data",
	// 	zap.Any("psql host", beconfig.GetConfig().Get("databases.postgres.host")),
	// 	zap.Any("psql port", beconfig.GetConfig().Get("databases.postgres.port")),
	// 	zap.Any("psql user", beconfig.GetConfig().Get("databases.postgres.user")),
	// 	zap.Any("psql password", beconfig.GetConfig().Get("databases.postgres.password")),
	// 	zap.Any("psql db", beconfig.GetConfig().Get("databases.postgres.db")),
	// 	zap.Any("psql sslmode", beconfig.GetConfig().Get("databases.postgres.sslmode")),
	// 	zap.Any("psql connect_timeout", beconfig.GetConfig().Get("databases.postgres.connect_timeout")),
	// 	zap.Any("mongo host", beconfig.GetConfig().Get("databases.mongo.host")),
	// 	zap.Any("mongo port", beconfig.GetConfig().Get("databases.mongo.port")),
	// 	zap.Any("mongo user", beconfig.GetConfig().Get("databases.mongo.user")),
	// 	zap.Any("mongo password", beconfig.GetConfig().Get("databases.mongo.password")),
	// 	zap.Any("mongo db", beconfig.GetConfig().Get("databases.mongo.db")),
	// 	zap.Any("mongo connect_timeout", beconfig.GetConfig().Get("databases.mongo.connect_timeout")),
	// )
	// //connect to psql db
	// psqlCfg := bedatabase.DbConfig{
	// 	Host:           beconfig.GetConfig().GetString("databases.postgres.host"),
	// 	Port:           beconfig.GetConfig().GetInt("databases.postgres.port"),
	// 	User:           beconfig.GetConfig().GetString("databases.postgres.user"),
	// 	Password:       beconfig.GetConfig().GetString("databases.postgres.password"),
	// 	Database:       beconfig.GetConfig().GetString("databases.postgres.db"),
	// 	SSLMode:        beconfig.GetConfig().GetString("databases.postgres.sslmode"),
	// 	ConnectTimeout: beconfig.GetConfig().GetInt("databases.postgres.connect_timeout"),
	// 	PostgresConfig: bedatabase.PostgresConfig{
	// 		Logger:           belogger.GetDBLogger(),
	// 		ApplicationName:  fmt.Sprintf("service-%s-%s-%d", os.Getenv("POD_NAME"), os.Getenv("POD_NAMESPACE"), os.Getpid()),
	// 		StatementTimeout: 5,  //TODO: make this configurable
	// 		MaxIdleConns:     10, //TODO: make this configurable
	// 		MaxOpenConns:     30, //TODO: make this configurable
	// 	},
	// }
	// postgresConn, err := bedatabase.ConnectPostgres(psqlCfg)
	// if err != nil {
	// 	belogger.Log().Error("Failed to connect psql database", zap.Error(err))
	// 	return err
	// }
	// databases = append(databases, postgresConn)

	// //connect to mongo db
	// mongoCfg := bedatabase.DbConfig{
	// 	Host:           beconfig.GetConfig().GetString("databases.mongo.host"),
	// 	Port:           beconfig.GetConfig().GetInt("databases.mongo.port"),
	// 	User:           beconfig.GetConfig().GetString("databases.mongo.user"),
	// 	Password:       beconfig.GetConfig().GetString("databases.mongo.password"),
	// 	Database:       beconfig.GetConfig().GetString("databases.mongo.db"),
	// 	ConnectTimeout: beconfig.GetConfig().GetInt("databases.mongo.connect_timeout"),
	// 	MongoConfig: bedatabase.MongoConfig{
	// 		ValidatePing: beconfig.GetConfig().GetBool("databases.mongo.ping"),
	// 	},
	// }
	// mClient, err := bedatabase.ConnectMongo(mongoCfg)
	// if err != nil {
	// 	belogger.Log().Error("Failed to connect mongo database", zap.Error(err))
	// 	return err
	// }
	// mongoDatabases = append(mongoDatabases, mClient)

	// //init middleware
	// //init repo, factory, service and controller interfaces
	// repo := dao.GetRepositoryItf(postgresConn, mClient)
	// service := service.GetServiceItf(repo)

	// // create and start the gin server
	// startRouter(service)

	return nil
}

// stops the router running in the go routine.
func ShutdownRouter() {
	// timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	// defer cancel()
	// belogger.Log().Info("Shutting down router START")
	// defer belogger.Log().Info("Shutting down router END")
	// if err := srv.Shutdown(timeoutCtx); err != nil {
	// 	belogger.Log().Fatal("Server forced to shutdown", zap.Error(err))
	// }
	// // catching ctx.Done(). timeout of 2 seconds.
	// <-timeoutCtx.Done()
	// log.Println("timeout of 2 seconds.")
}

// closes all database connections
//
//	closes each database connection made which are saved globally
//	logs error if unable to close
//		function used: *sql.DB.Close()
func CloseDatabase() {
	// belogger.Log().Info("disconnecting databases START")
	// defer belogger.Log().Info("disconnecting databases END")
	// for _, database := range databases {
	// 	db, _ := database.DB()
	// 	if db != nil {
	// 		err := db.Close()
	// 		if err != nil {
	// 			belogger.Log().Error("unable to close db", zap.Error(err))
	// 		}
	// 	} else {
	// 		belogger.Log().Error("unable to close db as connection is nil")
	// 	}
	// }
	// ctx := context.Background()
	// for _, database := range mongoDatabases {
	// 	database.Disconnect(ctx)
	// }
}

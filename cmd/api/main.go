package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"github.com/sauravkuila/mergemoney_assessment/pkg/constant"
	"github.com/sauravkuila/mergemoney_assessment/pkg/logger"
	"github.com/sauravkuila/mergemoney_assessment/pkg/server"
	"go.uber.org/zap/zapcore"
)

// package to start REST API server
func main() {
	var (
		environment string
	)
	host := os.Getenv("SERVER_HOST")
	if host != "" {
		environment = "server"
	} else {
		if len(os.Args) == 2 {
			environment = os.Args[1] // developer custom file
		} else {
			environment = "local"
		}
	}

	//initialize config
	config.Load(environment)

	//initialize logger
	logLevel, err := strconv.ParseInt(config.GetConfig().GetString("log.level"), 10, 64)
	if err != nil {
		log.Fatal("Invalid log config: ", err)
		os.Exit(1)
	}
	env := logger.DEVELOPMENT
	if config.GetConfig().GetString("env") != string(constant.DEVELOPMENT) {
		env = logger.PRODUCTION
	}
	if err := logger.InitLogger(env, zapcore.Level(logLevel)); err != nil {
		log.Fatal("Failed to initialize logger", err)
		os.Exit(1)
	}
	logger.Log().Info("Logger initialized")

	//start server with db connections
	if err := server.Start(); err != nil {
		log.Fatal("Failed to start server, err:", err)
		os.Exit(1)
	}
	addShutdownHook()
}

// adds and listens to any interrupt signal
//
//	the method should be called at the end of the main function as it blocks execution
//		internally it
//		-	closes redis connections
//		-	shuts down the http server
//		-	closes database connections
//	waits on syscall.SIGINT, syscall.SIGTERM, os.Interrupt
func addShutdownHook() {
	// when receive interruption from system shutdown server and scheduler
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	logger.Log().Info("Quit/Interrupt signal detected. Gracefully closing connections")
	//shutdown server
	server.ShutdownRouter()
	server.CloseDatabase()

	ctx := context.Background()

	logger.Log(ctx).Info(fmt.Sprintf("All done! Wrapping up here for PID: %d", os.Getpid()))
}

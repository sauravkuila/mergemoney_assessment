package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	config "github.com/sauravkuila/mergemoney_assessment/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logObject   *zap.Logger = nil
	dbLoggerObj *dbLoggerSt = nil
	err         error
)

const (
	DEVELOPMENT      = "development"
	PRODUCTION       = "production"
	syslogTimeFormat = "Jan 2 15:04:05"
)

func InitLogger(loggerType string, level zapcore.Level) error {
	var (
		cfg zap.Config
	)
	//switch based on logger type
	switch loggerType {
	case DEVELOPMENT:
		cfg = zap.NewDevelopmentConfig()
	case PRODUCTION:
		cfg = zap.NewProductionConfig()
	default:
		return fmt.Errorf("unsupported logger type")
	}

	//creating config for logger
	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.EncoderConfig.FunctionKey = "f"
	cfg.EncoderConfig.EncodeTime = syslogTimeEncoder
	cfg.EncoderConfig.ConsoleSeparator = " "
	cfg.Encoding = "console"

	//building logger
	logObject, err = cfg.Build()
	if err != nil {
		fmt.Println("failed to create custom production logger , Exiting system", err)
		os.Exit(0)
	} else if logObject == nil {
		logObject, err = getLoogerForType(loggerType)
		if err != nil {
			fmt.Println("failed to create production logger , Exiting system", err)
			os.Exit(0)
		}
		logObject.WithOptions(zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))
		fmt.Println("Failed to create custom production logger, creating production logger")
	} else {
		logObject.WithOptions(zap.AddCallerSkip(1), zap.AddStacktrace(zap.FatalLevel))
		fmt.Println("custom production logger created")
	}

	//building sugar logger
	sugaredLogger := logObject.Sugar()
	dbLoggerObj = &dbLoggerSt{sugaredLogger: sugaredLogger}

	return nil
}

func Log(data ...context.Context) *zap.Logger {
	if data != nil {
		ctx := data[0]
		return logObject.With(zap.Any(config.REQUESTID, ctx.Value(config.REQUESTID)), zap.Any(config.USERID, ctx.Value(config.USERID)), zap.Any(config.UCC, ctx.Value(config.UCC)))
	} else {
		return logObject
	}
}

func syslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(syslogTimeFormat))
}

func getLoogerForType(loggerType string) (*zap.Logger, error) {
	switch loggerType {
	case DEVELOPMENT:
		return zap.NewDevelopment()
	case PRODUCTION:
		return zap.NewProduction()
	default:
		return nil, fmt.Errorf("unsupported logger type")
	}
}

func GetDBLogger() *dbLoggerSt {
	return dbLoggerObj
}

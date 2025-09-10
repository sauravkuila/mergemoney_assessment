package logger

import (
	"context"
	"time"

	"go.uber.org/zap"
	gormlogger "gorm.io/gorm/logger"
)

type dbLoggerSt struct {
	sugaredLogger *zap.SugaredLogger
}

func (obj *dbLoggerSt) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return &dbLoggerSt{sugaredLogger: obj.sugaredLogger}
}

func (obj *dbLoggerSt) Info(ctx context.Context, msg string, args ...interface{}) {
	obj.sugaredLogger.With(ctx).Infof(msg, args...)
}

func (obj *dbLoggerSt) Error(ctx context.Context, message string, args ...interface{}) {
	obj.sugaredLogger.With(ctx).Errorf(message, args...)
}

func (obj *dbLoggerSt) Warn(ctx context.Context, message string, args ...interface{}) {
	obj.sugaredLogger.With(ctx).Warnf(message, args...)
}

func (obj *dbLoggerSt) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		obj.sugaredLogger.Errorw("gorm trace", "error", err, "elapsed", elapsed, "sql", sql, "rows", rows)
	} else {
		obj.sugaredLogger.Infow("gorm trace", "elapsed", elapsed, "sql", sql, "rows", rows)
	}
}

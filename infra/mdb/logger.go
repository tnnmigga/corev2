package mdb

import (
	"context"
	"time"

	"github.com/tnnmigga/corev2/log"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

type gormLogger struct{}

func (l gormLogger) LogMode(logger.LogLevel) logger.Interface {
	return l
}

func (l gormLogger) Info(ctx context.Context, f string, s ...interface{}) {
	log.Infof(f, s...)
}

func (l gormLogger) Warn(ctx context.Context, f string, s ...interface{}) {
	log.Warnf(f, s...)
}

func (l gormLogger) Error(ctx context.Context, f string, s ...interface{}) {
	log.Errorf(f, s...)
}

func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if err != nil {
		sql, ra := fc()
		log.Errorf("exec sql error %v, SQL: %s, rows affected: %d", err, sql, ra)
	} else if log.Logger().Level() == zapcore.DebugLevel {
		sql, ra := fc()
		log.Debugf("exec SQL: %s, rows affected: %d, time cost: %v", sql, ra, time.Since(begin))
	}
}

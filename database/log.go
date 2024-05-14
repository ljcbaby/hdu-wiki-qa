package database

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type LogrusLogger struct {
	logger *logrus.Logger
}

func (l *LogrusLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *LogrusLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithField("module", "gorm").WithContext(ctx).Infof(msg, data...)
}

func (l *LogrusLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithField("module", "gorm").WithContext(ctx).Warnf(msg, data...)
}

func (l *LogrusLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithField("module", "gorm").WithContext(ctx).Errorf(msg, data...)
}

func (l *LogrusLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.logger.WithField("module", "gorm").WithContext(ctx).WithField("duration", time.Since(begin).Seconds()).Debugf(fc())
}

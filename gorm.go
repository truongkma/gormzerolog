package gormzerolog

import (
	"context"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	gormLogger "gorm.io/gorm/logger"
)

type Logger struct {
	logger zerolog.Logger
	config Config
}

type Config struct {
	SlowThreshold        time.Duration
	ParameterizedQueries bool
}

func NewLogger(config Config) *Logger {
	return &Logger{
		logger: log.Logger,
		config: config,
	}
}

func (l *Logger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	newLogger := l.logger.Level(zerolog.InfoLevel)
	switch level {
	case gormLogger.Silent:
		newLogger = l.logger.Level(zerolog.NoLevel)
	case gormLogger.Error:
		newLogger = l.logger.Level(zerolog.ErrorLevel)
	case gormLogger.Warn:
		newLogger = l.logger.Level(zerolog.WarnLevel)
	case gormLogger.Info:
		newLogger = l.logger.Level(zerolog.InfoLevel)
	}
	return &Logger{logger: newLogger, config: l.config}
}

func (l *Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info().Ctx(ctx).Msgf(msg, data...)
}

func (l *Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn().Ctx(ctx).Msgf(msg, data...)
}

func (l *Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error().Ctx(ctx).Msgf(msg, data...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	sqlStr := strings.ReplaceAll(sql, "\"", "")

	if err != nil {
		l.logger.Error().Ctx(ctx).CallerSkipFrame(3).Dur("elapsed", elapsed).Str("sql", sqlStr).Int64("rows", rows).Err(err).Msg("Error SQL")
	} else {
		if l.config.SlowThreshold != 0 && elapsed >= l.config.SlowThreshold {
			l.logger.Warn().Ctx(ctx).CallerSkipFrame(3).Dur("elapsed", elapsed).Str("sql", sqlStr).Int64("rows", rows).Msg("Warn SQL")
		} else {
			if l.logger.GetLevel() != zerolog.TraceLevel && zerolog.GlobalLevel() > zerolog.InfoLevel {
				l.logger.Log().Ctx(ctx).CallerSkipFrame(3).Dur("elapsed", elapsed).Str("sql", sqlStr).Int64("rows", rows).Msg("Debug SQL")
			} else {
				l.logger.Info().Ctx(ctx).CallerSkipFrame(3).Dur("elapsed", elapsed).Str("sql", sqlStr).Int64("rows", rows).Msg("Trace SQL")
			}
		}
	}
}

func (l *Logger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.config.ParameterizedQueries {
		return sql, nil
	}

	return sql, params
}

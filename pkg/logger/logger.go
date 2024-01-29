package logger

import (
	"boilerplate-service/constant"
	"context"

	"go.uber.org/zap"
)

type ILogger interface {
	Debug(ctx context.Context, msg string, fields ...zap.Field)
	Info(ctx context.Context, msg string, fields ...zap.Field)
	Error(ctx context.Context, msg string, fields ...zap.Field)
	Warn(ctx context.Context, msg string, fields ...zap.Field)
	Panic(ctx context.Context, msg string, fields ...zap.Field)
	Sync() error

	GetLogger() *zap.Logger
}

type Config struct {
	Environment string
	ServiceName string
}

type logger struct {
	zapLog *zap.Logger
}

func New(config Config) (ILogger, error) {
	zapConfig := zap.Config{}

	if config.Environment == constant.EnvironmentLocal || config.Environment == constant.EnvironmentDevelopment {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zapConfig.Development = true
	} else {
		zapConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
		zapConfig.Development = false
		zapConfig.Sampling = &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		}
	}

	zapConfig.Encoding = "json"
	zapConfig.EncoderConfig = zap.NewProductionEncoderConfig()
	zapConfig.OutputPaths = []string{"stdout"}
	zapConfig.ErrorOutputPaths = []string{"stderr"}
	zapConfig.InitialFields = map[string]interface{}{
		"env": config.Environment,
		"app": config.ServiceName,
	}

	zapLog, err := zapConfig.Build()
	return &logger{zapLog}, err
}

func (l *logger) getTraceId(ctx context.Context) string {
	if ctx.Value(constant.CtxTraceIdKey) != nil {
		return ctx.Value(constant.CtxTraceIdKey).(string)
	}
	return ""
}

func (l *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	trace := l.getTraceId(ctx)
	if trace != "" {
		fields = append(fields, zap.String("trace_id", trace))
	}
	l.zapLog.Debug(msg, fields...)
}

func (l *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	trace := l.getTraceId(ctx)
	if trace != "" {
		fields = append(fields, zap.String("trace_id", trace))
	}
	l.zapLog.Info(msg, fields...)
}

func (l *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	trace := l.getTraceId(ctx)
	if trace != "" {
		fields = append(fields, zap.String("trace_id", trace))
	}
	l.zapLog.Error(msg, fields...)
}

func (l *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	trace := l.getTraceId(ctx)
	if trace != "" {
		fields = append(fields, zap.String("trace_id", trace))
	}
	l.zapLog.Warn(msg, fields...)
}

func (l *logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	trace := l.getTraceId(ctx)
	if trace != "" {
		fields = append(fields, zap.String("trace_id", trace))
	}
	l.zapLog.Panic(msg, fields...)
}

func (l *logger) Sync() error {
	return l.zapLog.Sync()
}

func (l *logger) GetLogger() *zap.Logger {
	return l.zapLog
}

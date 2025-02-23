package logger

import (
	"errors"
	"syscall"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.SugaredLogger, func(), error) {
	cfg := &zap.Config{
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		Sampling: nil,
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "name",
			CallerKey:      "caller",
			FunctionKey:    "func",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339),
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	zapLogger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, nil, err
	}

	flushFn := func() {
		if err = zapLogger.Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
			zapLogger.Error("unable to flush logs", zap.Error(err))
		}
	}

	return zapLogger.Sugar(), flushFn, nil
}

package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

// NewZap - initialize new logger.
func NewZap(level string) Logger {
	var l zapcore.Level
	switch strings.ToLower(level) {
	case "error":
		l = zapcore.ErrorLevel
	case "warn":
		l = zapcore.WarnLevel
	case "info":
		l = zapcore.InfoLevel
	case "debug":
		l = zapcore.DebugLevel
	default:
		l = zapcore.InfoLevel
	}

	// Create logger config
	config := zap.Config{
		Development:      false,
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(l),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			EncodeDuration: zapcore.SecondsDurationEncoder,
			LevelKey:       "severity",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			TimeKey:        "timestamp",
			EncodeTime:     zapcore.ISO8601TimeEncoder, // e.g. 2020-05-05T03:24:36.903+0300
			NameKey:        "name",
			EncodeName:     zapcore.FullNameEncoder, // e.g. GetSiteGeneralHandler
			MessageKey:     "message",
			EncodeCaller:   zapcore.ShortCallerEncoder, // e.g. package/file:line
			CallerKey:      "caller",
			StacktraceKey:  "",
			LineEnding:     "\n",
		},
	}

	// Build logger from config
	logger, _ := config.Build()
	return Logger{
		logger.Sugar(),
	}
}

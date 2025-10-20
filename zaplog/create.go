package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogName string

const (
	LogNameLogger   LogName = "logger"
	LogNameEcho     LogName = "echo"
	LogNameReductor LogName = "reductor"
	LogNameTrue     LogName = "true"
)

type LogConfig struct {
	ErrorOutputPaths []string
	Debug            bool
	Console          bool // true - "console" false - "json"
	Name             string
}

func isValidLogName(s string) bool {
	switch LogName(s) {
	case LogNameLogger, LogNameEcho, LogNameReductor, LogNameTrue:
		return true
	default:
		return false
	}
}

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "ts",
	LevelKey:       "lvl",
	NameKey:        "logger",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

// createLogger creates a zap logger with specified output paths and debug level.
// console determines the encoding format: true for console, false for JSON.
func createLogger(output []string, debug bool, console bool) (*zap.Logger, error) {
	level := zap.InfoLevel
	if debug {
		level = zap.DebugLevel
	}
	if len(output) == 0 {
		output = []string{"stderr"}
	}

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(level),
		Development:       debug,
		DisableCaller:     false,
		DisableStacktrace: !debug,
		Encoding:          "console",
		Sampling:          nil,
		EncoderConfig:     encoderConfig,
		OutputPaths:       output,
		ErrorOutputPaths:  []string{"stderr"},
	}
	if !console {
		config.Encoding = "json"
	}

	if !debug {
		config.Sampling = &zap.SamplingConfig{Initial: 100, Thereafter: 100}
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return logger, nil
}

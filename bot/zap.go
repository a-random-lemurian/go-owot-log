package bot

import (
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetZapLogger(v *viper.Viper, filename string) *zap.Logger {
	// Define level enablers
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel 
	})

	goToConsole := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	// Create encoders
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "stacktrace",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "stacktrace",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// Configure rolling log file
	logFile := zapcore.AddSync(&lumberjack.Logger{
		Filename:   v.GetString("data_directory") + "/" + filename,
		MaxSize:    200, // megabytes
		Compress:   true,
	})

	// Create cores
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.Lock(os.Stderr),
		goToConsole,
	)

	fileCore := zapcore.NewCore(
		fileEncoder,
		logFile,
		highPriority,
	)

	// Combine cores
	core := zapcore.NewTee(
		consoleCore,
		fileCore,
	)

	// Create logger with options
	log := zap.New(
		core,
		zap.AddCaller(),                       // Adds caller information
	)

	return log
}

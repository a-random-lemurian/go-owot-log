package bot

import (
	"os"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var fileEncoder = zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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

func GetZapLogger(v *viper.Viper, filename string) *zap.Logger {
	var l zapcore.Level
	unmarshalErr := l.UnmarshalText([]byte(v.GetString("log.minimum_level")))
	if unmarshalErr != nil {
		l = zapcore.InfoLevel
	}

	goToConsole := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= l
	})

	// Create cores
	consoleCore := zapcore.NewCore(
		fileEncoder,
		zapcore.Lock(os.Stdout),
		goToConsole,
	)

	// Create logger with options
	log := zap.New(
		consoleCore,
		zap.AddCaller(), // Adds caller information
	)

	if unmarshalErr != nil {
		log.Error("Failed to determine log level from file", zap.Error(unmarshalErr))
	}

	return log
}

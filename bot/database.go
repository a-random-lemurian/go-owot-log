package bot

import (
	"fmt"

	"github.com/a-random-lemurian/go-owot-log/internal/types"
	"github.com/natefinch/lumberjack"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Database interface {
	AddMessage(*types.Message) error
}

type PlaintextFileDatabase struct {
	Database
	owotLog *zap.Logger
}

func getPlaintextZap(v *viper.Viper, filename string) *zap.Logger {
	logFile := zapcore.AddSync(&lumberjack.Logger{
		Filename: v.GetString("data_directory") + "/" + filename,
		MaxSize:  200, // megabytes
		Compress: true,
	})
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})
	fileCore := zapcore.NewCore(
		fileEncoder,
		logFile,
		highPriority,
	)
	log := zap.New(fileCore, zap.AddCaller())
	return log
}

func CreatePlainTextFileDatabase(v *viper.Viper, parentLog *zap.Logger) (*PlaintextFileDatabase, error) {
	log := parentLog.Named("plaintextdb")
	file := "owot.log"
	log.Info(fmt.Sprintf("Creating plaintext file database at %s", file), zap.String("filename", file))
	db := PlaintextFileDatabase{}
	db.owotLog = getPlaintextZap(v, file).Named("zapdb")
	return &db, nil
}

func (pfd *PlaintextFileDatabase) AddMessage(message *types.Message) error {
	pfd.owotLog.Info("got message",
		zap.ByteString("rawMessage", message.WebsocketMessage),
		zap.String("originWebsocketUrl", message.Origin),
		zap.Time("received", message.Time),
	)
	return nil
}

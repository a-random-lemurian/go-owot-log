package bot

import (
	"github.com/a-random-lemurian/go-owot-log/internal/types"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Database interface {
	AddMessage(*types.Message) error
}

type PlaintextFileDatabase struct {
	Database
	owotLog *zap.Logger
}

func CreatePlainTextFileDatabase(v *viper.Viper) (*PlaintextFileDatabase, error) {
	db := PlaintextFileDatabase{}
	db.owotLog = GetZapLogger(v, "owot.log")
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

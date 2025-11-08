package logger

import (
	"github.com/a-random-lemurian/go-owot-log/config"
	"go.uber.org/zap"
)

type Database interface {
	AddMessage(*Message) error
}

type PlaintextFileDatabase struct {
	Database
	owotLog *zap.Logger
}

func CreatePlainTextFileDatabase(cfg *config.Config) (*PlaintextFileDatabase, error) {
	db := PlaintextFileDatabase{}
	db.owotLog = GetZapLogger(cfg, "owot.log")
	return &db, nil
}

func (pfd *PlaintextFileDatabase) AddMessage(message *Message) error {
	pfd.owotLog.Info("got message",
		zap.ByteString("rawMessage", message.WebsocketMessage),
		zap.String("originWebsocketUrl", message.Origin),
		zap.Time("received", message.Time),
	)
	return nil
}

package bot

import (
	"fmt"

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

func CreatePlainTextFileDatabase(v *viper.Viper, parentLog *zap.Logger) (*PlaintextFileDatabase, error) {
	log := parentLog.Named("plaintextdb")
	file := "owot.log"
	log.Info(fmt.Sprintf("Creating plaintext file database at %s", file), zap.String("filename", file))
	db := PlaintextFileDatabase{}
	db.owotLog = GetZapLogger(v, file).Named("zapdb")
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

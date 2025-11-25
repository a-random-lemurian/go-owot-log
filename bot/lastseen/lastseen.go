package lastseen

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/a-random-lemurian/go-owot-log/bot/chatcmd"

	"github.com/a-random-lemurian/go-owot"

	slogGorm "github.com/orandin/slog-gorm"
	"github.com/pkg/errors"
	slogzap "github.com/samber/slog-zap/v2"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Lastseen struct {
	db  LastseenDB
	log *zap.Logger
}

// getLsDB: get lastseen database
func getLsDB(v *viper.Viper, log *zap.Logger) (LastseenDB, error) {
	if v.GetString("database.type") != "sqlite" {
		return nil, errors.New("unsupported database")
	}

	slogHandler := slogzap.Option{
		Level:  slog.LevelDebug,
		Logger: log.Named("database"),
	}.NewZapHandler()

	gormLogger := slogGorm.New(
		slogGorm.WithHandler(slogHandler), // pass the handler
	)

	sqlite, err := gorm.Open(sqlite.Open(v.GetString("database.sqlite.dsn")), &gorm.Config{Logger: gormLogger})
	if err != nil {
		return nil, errors.Wrap(err, "initializing database")
	}
	db := NewGormLastseenDB(sqlite)
	return db, nil
}

func NewLastseen(parentV *viper.Viper, log *zap.Logger) (*Lastseen, error) {
	v := parentV.Sub("lastseen")

	sublog := log.Named("lastseen")
	ls := &Lastseen{
		log: sublog,
	}

	db, err := getLsDB(v, ls.log)
	if err != nil {
		return nil, errors.Wrap(err, "initializing lastseen")
	}
	ls.db = db

	sublog.Info("Initialized lastseen")

	return ls, nil
}

// Ingest a new chat message to update the last seen times.
func (l *Lastseen) IngestMessage(chat *owot.MessageChat) {
	l.db.SetLastseen(chat.RealUsername, chat.Date.Time())
}

// Consent should be "yes" or "no"
func (l *Lastseen) validateConsent(s string) (bool, error) {
	switch s {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, errors.New("invalid consent setting, use yes or no")
	}
}

func (l *Lastseen) cmdLastseen(chatCtx chatcmd.Context, c *owot.MessageChat) {
	tokens := strings.Fields(c.Message)
	if len(tokens) < 3 {
		chatCtx.Chat("usage: ch lastseen <user>")
		return
	}

	target := tokens[2]
	consent := l.db.Consents(target)
	if !consent {
		chatCtx.Chat(fmt.Sprintf("%s has opted out", target))
		return
	}
	lastseen, err := l.db.GetLastseen(target)
	lastseenTimeString := lastseen.UTC().Format("2006-01-02 15:04:05.000") + " UTC"
	if err != nil {
		chatCtx.Chat(fmt.Sprintf("%s not found", target))
		return
	}
	chatCtx.Chat(fmt.Sprintf("%s was last seen at %s", target, lastseenTimeString))
}

func (l *Lastseen) cmdConsent(chatCtx chatcmd.Context, c *owot.MessageChat) {
	tokens := strings.Fields(c.Message)

	if len(tokens) < 3 {
		chatCtx.Chat("usage: ch lastseen-consent {yes|no}")
		return
	}

	consent, err := l.validateConsent(tokens[2])
	if err != nil {
		chatCtx.Chat(fmt.Sprintf("error: %s", err))
		return
	}

	// 1. user opting in and out
	if len(tokens) < 4 {
		err := l.db.SetConsent(c.RealUsername, consent)
		if err != nil {
			l.log.Error("database error while setting consent", zap.Error(err))
			chatCtx.Chat("there was a database error. sorry")
			return
		}

		if consent {
			chatCtx.Chat("You have opted back in.")
		} else {
			chatCtx.Chat("You have opted out.")
		}

		l.log.Info("user changed consent setting",
			zap.String("realUsername", c.RealUsername),
			zap.Bool("consent", consent),
		)

		return
	}

	// 2. forcefully optin/out someone (admin)
	if !chatCtx.IsTrusted(c.RealUsername) {
		return
	}

	if len(tokens) < 4 {
		chatCtx.Chat("user to change consent for unspecified")
		return
	}

	l.adminChangeConsent(chatCtx, consent, tokens[3])
}

func (l *Lastseen) adminChangeConsent(chatCtx chatcmd.Context, consent bool, target string) error {
	err := l.db.SetConsent(target, consent)
	if err != nil {
		chatCtx.Chat("there was a database error. sorry")
		return err
	}
	chatCtx.Chat(fmt.Sprintf("changed optin/out of %s to %t", target, consent))
	return nil
}

func (l *Lastseen) AttachCommands(b chatcmd.CommandRegistry) error {
	b.RegisterCommand("lastseen", l.cmdLastseen)
	b.RegisterCommand("lastseen-consent", l.cmdConsent)
	return nil
}

package bot

import (
	"sync"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Bot struct {
	v  *viper.Viper
	log *zap.Logger
	db Database
}

func NewBot(v *viper.Viper) (*Bot, error) {
	return &Bot{v: v}, nil
}

func (b *Bot) getDatabase() error {
	database, err := CreatePlainTextFileDatabase(b.v)
	if err != nil {
		b.log.Error("Cannot create PlainTextFileDatabase", zap.Error(err))
		return err
	}
	b.db = database
	return nil
}

func (b *Bot) Run() error {
	log := GetZapLogger(b.v, "application.log")
	b.log = log

	if err := b.getDatabase(); err != nil {
		return err
	}

	log.Info("Hi, good to see you")
	defer log.Info("Goodbye")

	wg := sync.WaitGroup{}

	worlds := b.v.GetStringSlice("worlds")
	globalReceiver := b.v.GetString("global_receiver")

	for _, worldName := range worlds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Creating world", zap.String("worldName", worldName))
			acceptsGlobal := (worldName == globalReceiver)
			for {
				world, err := CreateWorld(worldName, log, b.db, acceptsGlobal)
				if err != nil {
					log.Error("Cannot create World", zap.String("worldName", worldName), zap.Error(err))
					time.Sleep(1 * time.Second)
					continue
				}
				err = world.Run()
				if err != nil {
					log.Warn("Disconnected from world, reconnecting", zap.String("worldName", worldName), zap.Error(err))
				}
				time.Sleep(1 * time.Second)
			}
		}()
	}

	wg.Wait()

	return nil
}
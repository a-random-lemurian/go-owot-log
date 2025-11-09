package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/a-random-lemurian/go-owot-log/bot"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

func main() {
	app := &cli.App{
		Usage:   "A program for monitoring OWOT chat",
		Authors: []*cli.Author{{Name: "Lemuria"}},
		Commands: []*cli.Command{
			{
				Name:   "start",
				Usage:  "Start the logger",
				Action: cmdStart,
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}
}

func cmdStart(cCtx *cli.Context) error {
	v := viper.New()
	v.AddConfigPath(".")

	v.SetEnvPrefix("OWOTLOG")
	v.AutomaticEnv()

	v.SetConfigName("config")
	v.SetConfigType("yml")

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	mandatoryKeys := []string{"worlds", "data_directory", "global_receiver"}
	for _, key := range mandatoryKeys {
		if !v.IsSet(key) {
			fmt.Printf("your settings: %v\n", v.AllSettings())
			return fmt.Errorf("key %s is required but not set in config", key)
		}
	}

	log := bot.GetZapLogger(v, "application.log")
	log.Info("Hi, good to see you")
	defer log.Info("Goodbye")

	database, err := bot.CreatePlainTextFileDatabase(v)
	if err != nil {
		log.Error("Cannot create PlainTextFileDatabase", zap.Error(err))
		return err
	}

	wg := sync.WaitGroup{}

	worlds := v.GetStringSlice("worlds")
	globalReceiver := v.GetString("global_receiver")

	for _, worldName := range worlds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Creating world", zap.String("worldName", worldName))
			acceptsGlobal := (worldName == globalReceiver)
			for {
				world, err := bot.CreateWorld(worldName, log, database, acceptsGlobal)
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

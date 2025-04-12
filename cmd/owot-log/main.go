package main

import (
	"os"
	"sync"
	"time"

	"github.com/a-random-lemurian/go-owot-log/config"
	"github.com/a-random-lemurian/go-owot-log/logger"
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

	app.Run(os.Args)
}

func cmdStart(cCtx *cli.Context) error {
	cfg, err := config.LoadDefault()
	if err != nil {
		return err
	}
	log := logger.GetZapLogger(cfg, "application.log")
	log.Info("Hi, good to see you")
	defer log.Info("Goodbye")

	database, err := logger.CreatePlainTextFileDatabase(cfg)
	if err != nil {
		log.Error("Cannot create PlainTextFileDatabase", zap.Error(err))
		return err
	}

	wg := sync.WaitGroup{}

	for _, worldName := range cfg.Worlds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Info("Creating world", zap.String("worldName", worldName))
			acceptsGlobal := (worldName == cfg.GlobalReceiver)
			for {
				world, err := logger.CreateWorld(worldName, log, database, acceptsGlobal)
				if err != nil {
					log.Error("Cannot create World", zap.String("worldName", worldName), zap.Error(err))
					time.Sleep(1 * time.Second)
				}
				world.Run()
			}
		}()
	}

	wg.Wait()

	return nil
}

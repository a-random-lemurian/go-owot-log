package main

import (
	"fmt"
	"os"

	"github.com/a-random-lemurian/go-owot-log/bot"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
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

	bot, err := bot.NewBot(v)
	if err != nil {
		return err
	}

	bot.Run()

	return nil
}

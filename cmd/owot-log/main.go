package main

import (
	"fmt"
	"strings"

	"github.com/a-random-lemurian/go-owot-log/bot"
	"github.com/carlmjohnson/versioninfo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cmdRoot = &cobra.Command{
	Use:     "owot-log",
	Short:   "A bot for OWOT chat",
	Version: versioninfo.Short(),
}

var cmdStart = &cobra.Command{
	Use:   "start",
	Short: "Start the bot",
	RunE:  runStart,
}

func init() {
	cmdRoot.AddCommand(cmdStart)

	csF := cmdStart.Flags()
	csF.String("config", "config.yml", "Path to configuration file.")
}

func main() {
	cmdRoot.Execute()
}

func runStart(cmd *cobra.Command, args []string) error {
	configPath, err := cmd.Flags().GetString("config")
	if err != nil {
		return err
	}

	v := viper.New()
	v.AddConfigPath(".")

	v.SetEnvPrefix("OWOTLOG")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.SetConfigFile(configPath)

	err = v.ReadInConfig()
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

	err = bot.Run()

	return err
}

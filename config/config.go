package config

import (
	"encoding/json"
	"os"

	"github.com/natefinch/lumberjack"
)

type Config struct {
	Lumberjack     lumberjack.Logger `json:"lumberjack"`
	DataDirectory  string            `json:"dataDirectory"`
	Worlds         []string          `json:"worlds"`
	GlobalReceiver string            `json:"globalReceiver"`
}

// Load a configuration from a []byte. This also calls GetChannels() to read the
// channel file.
func ParseConfig(jsonData []byte) (*Config, error) {
	var config Config

	err := json.Unmarshal(jsonData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Load a configuration from a file.
func ConfigFromFile(filename string) (*Config, error) {
	file, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return ParseConfig(file)
}

// Load a configuration from the default file, config.json.
func LoadDefault() (*Config, error) {
	return ConfigFromFile("config.json")
}

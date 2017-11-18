package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	MessengerToken string `json:"messenger_token"`
	LogsLinesCount int    `json:"log_lines_count"`
}

func FromFile(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = json.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

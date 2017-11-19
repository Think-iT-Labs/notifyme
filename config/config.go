package config

import (
	"io/ioutil"
	"os"

	"github.com/yosuke-furukawa/json5/encoding/json5"
)

type Config struct {
	MessengerEnabled bool     `json:"messenger_enabled"`
	MessengerTokens  []string `json:"messenger_tokens"`
}

func FromFile(filename string) (Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return Config{}, err
	}

	var config Config

	err = json5.NewDecoder(file).Decode(&config)
	if err != nil {
		return Config{}, err
	}

	return config, nil
}

func CreateDefault(path string) error {
	configTemplate := `{
		// This option control whenever messenger notifications are enabled or no
		"messenger_enabled": true,
		
		// Append your messenger tokens here, 
		// If you don't now you token, ask the Facebook Chat Bot for it.
		// You can talk to the Chat Bot by sending a message to: https://www.facebook.com/clinotify.me
		"messenger_tokens": [
			""
		],
	
		// This option control when notifications should be send.  
		// Should be one of: "all", "success_only" or "error_only"
		"enable_for_status": "all"
	}
	`
	return ioutil.WriteFile(path, []byte(configTemplate), 0644)
}

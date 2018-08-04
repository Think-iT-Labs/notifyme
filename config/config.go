package config

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

var DefaultConfigPath string

func init() {
	user, _ := user.Current()
	DefaultConfigPath = filepath.Join(user.HomeDir, ".notifyme")
}

type Config struct {
	Carriers       []map[string]interface{} `yaml:"carriers"`
	withTimestamps bool                     `yaml:"withTimestamps"`
}

// FromFile load a configuration given a filename
func FromFile(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	config := new(Config)

	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// CreateDefault creates a sample configutation file at the path `DefaultConfigPath`
func CreateDefault() error {
	configTemplate := `
---
withTimestamps: false

carriers:
  - type: slack
    token: "xoxp-XXXXXX"
    channels: "@user1, #general"
`
	return ioutil.WriteFile(DefaultConfigPath, []byte(configTemplate), 0644)
}

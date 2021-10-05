package config

import (
	"os"
	"encoding/json"
)

type Config struct {
	Twitter TwitterConfig `json:"twitter"`
}

type TwitterConfig struct {
	ConsumerKey string `json:"consumer_key"`
	ConsumerSecret string `json:"consumer_secret"`
	AccessToken string `json:"access_key"`
	AccessSecret string `json:"access_secret"`
}

func LoadConfig(file string) (config *Config, err error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	parser := json.NewDecoder(f)
	parser.Decode(&config)
	return
}

package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/caarlos0/env/v6"
)

// Config root of config
type Config struct {
	UserName     string             `json:"userName" env:"JOY_USER_NAME"`
	TwitterCreds TwitterCredsConfig `json:"twitterCreds"`
}

// TwitterCredsConfig creds for Twitter API access
type TwitterCredsConfig struct {
	AccessToken       string `json:"accessToken" env:"TW_ACCESS_TOKEN"`
	AccessTokenSecret string `json:"accessTokenSecret" env:"TW_ACCESS_TOKEN_SECRET"`
	ConsumerKey       string `json:"consumerKey" env:"TW_CONSUMER_KEY"`
	ConsumerSecret    string `json:"consumerSecret" env:"TW_CONSUMER_SECRET"`
}

// ReadFromJSONFile reads and deserialize config from json file
func ReadFromJSONFile(path string) (*Config, error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Error reading file: %s", err.Error())
	}
	var config Config
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, fmt.Errorf("Error parsing json file: %s", err.Error())
	}
	return &config, nil
}

// ReadFromEnv reads config from env
func ReadFromEnv() (*Config, error) {
	cfg := Config{}
	err := env.Parse(&cfg)
	return &cfg, err
}

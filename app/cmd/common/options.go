package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"joytotwi/app/twisender"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v6"
)

// CommonOptions for all commands
type Options struct {
	// see ParserID consts in each of parsers impls
	SourceType        string `json:"sourceType" env:"JOY_SRC_TYPE" envDefault:"page"`
	UserName          string `json:"userName" env:"JOY_USER_NAME"`
	AccessToken       string `json:"accessToken" env:"TW_ACCESS_TOKEN"`
	AccessTokenSecret string `json:"accessTokenSecret" env:"TW_ACCESS_TOKEN_SECRET"`
	ConsumerKey       string `json:"consumerKey" env:"TW_CONSUMER_KEY"`
	ConsumerSecret    string `json:"consumerSecret" env:"TW_CONSUMER_SECRET"`
}

// Commander is command containing CommonOptions
type Commander interface {
	SetCommonOptions(opts *Options)
}

// ReadFromJSONFile reads options from json config
func (opts *Options) ReadFromJSONFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid file path: '%s', %s", path, err.Error())
	}
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("error reading file: %s", err.Error())
	}
	err = json.Unmarshal(bytes, opts)
	if err != nil {
		return fmt.Errorf("error parsing json file: %s", err.Error())
	}
	return nil
}

// ReadFromEnv reads options from env
func (opts *Options) ReadFromEnv() error {
	return env.Parse(opts)
}

// Validate option values
func (opts *Options) Validate() error {
	var messages []string

	checkNotEmpty := func(val string, name string) {
		if strings.Trim(val, " ") == "" {
			messages = append(messages, name+" is empty")
		}
	}

	checkNotEmpty(opts.SourceType, "source type")
	checkNotEmpty(opts.UserName, "user name")
	checkNotEmpty(opts.AccessToken, "access token")
	checkNotEmpty(opts.AccessTokenSecret, "access token secret")
	checkNotEmpty(opts.ConsumerKey, "consumer key")
	checkNotEmpty(opts.ConsumerSecret, "consumer secret")

	if len(messages) > 0 {
		return errors.New(strings.Join(messages, ", "))
	}
	return nil
}

// GetTwiCreds extracts twi creds
func (opts Options) GetTwiCreds() twisender.ClientCreds {
	return twisender.ClientCreds{
		AccessToken:       opts.AccessToken,
		AccessTokenSecret: opts.AccessTokenSecret,
		ConsumerKey:       opts.ConsumerKey,
		ConsumerSecret:    opts.ConsumerSecret,
	}
}

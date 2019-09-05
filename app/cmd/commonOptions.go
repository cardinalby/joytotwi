package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v6"
)

// CommonOptions for all commands
type CommonOptions struct {
	// see ParserID consts in each of parsers impls
	SourceType        string `json:"sourceType" env:"JOY_SRC_TYPE" envDefault:"page"`
	UserName          string `json:"userName" env:"JOY_USER_NAME"`
	AccessToken       string `json:"accessToken" env:"TW_ACCESS_TOKEN"`
	AccessTokenSecret string `json:"accessTokenSecret" env:"TW_ACCESS_TOKEN_SECRET"`
	ConsumerKey       string `json:"consumerKey" env:"TW_CONSUMER_KEY"`
	ConsumerSecret    string `json:"consumerSecret" env:"TW_CONSUMER_SECRET"`
}

// CommonOptionsCommander is command containing CommonOptions
type CommonOptionsCommander interface {
	SetCommonOptions(opts *CommonOptions)
}

// SetFromAppOptions assign correspondent fields from appOpts
func (opts *CommonOptions) SetFromAppOptions(appOpts *AppOptions) {
	opts.SourceType = appOpts.SourceType
	opts.UserName = appOpts.UserName
	opts.AccessToken = appOpts.AccessToken
	opts.AccessTokenSecret = appOpts.AccessTokenSecret
	opts.ConsumerKey = appOpts.ConsumerKey
	opts.ConsumerSecret = appOpts.ConsumerSecret
}

// ReadFromJSONFile reads options from json config
func (opts *CommonOptions) ReadFromJSONFile(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("Invalid file path: '%s', %s", path, err.Error())
	}
	bytes, err := ioutil.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("Error reading file: %s", err.Error())
	}
	err = json.Unmarshal(bytes, opts)
	if err != nil {
		return fmt.Errorf("Error parsing json file: %s", err.Error())
	}
	return nil
}

// ReadFromEnv reads options from env
func (opts *CommonOptions) ReadFromEnv() error {
	return env.Parse(opts)
}

// Validate option values
func (opts *CommonOptions) Validate() error {
	messages := []string{}

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

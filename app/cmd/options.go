package cmd

import (
	"joytotwi/app/cmd/common"
	"joytotwi/app/cmd/fill"
	"joytotwi/app/cmd/update"
	"joytotwi/app/cmd/watch"
)

// AppOptions with all cli commands and flags
type AppOptions struct {
	Update update.Command `command:"update"`
	Watch  watch.Command  `command:"watch"`
	Fill   fill.Command   `command:"fill"`

	ConfigFile string `short:"c" long:"config-file" default:"" description:"Load app options from json file"`

	SourceType        string `short:"s" long:"source-type" default:"page" description:"Posts source: rss or page"`
	UserName          string `short:"u" long:"user-name" description:"JoyReactor user name"`
	AccessToken       string `long:"access-token" description:"Twitter API AccessToken"`
	AccessTokenSecret string `long:"access-token-secret" description:"Twitter API Access Token Secret"`
	ConsumerKey       string `long:"consumer-key" description:"Twitter API Consumer Key"`
	ConsumerSecret    string `long:"consumer-secret" description:"Twitter API Consumer Secret"`
}

// GetCommonOptions returns subset of options needed for every command
func (opts *AppOptions) GetCommonOptions() *common.Options {
	return &common.Options{
		SourceType:        opts.SourceType,
		UserName:          opts.UserName,
		AccessToken:       opts.AccessToken,
		AccessTokenSecret: opts.AccessTokenSecret,
		ConsumerKey:       opts.ConsumerKey,
		ConsumerSecret:    opts.ConsumerSecret,
	}
}

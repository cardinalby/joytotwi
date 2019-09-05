package cmd

import (
	log "github.com/sirupsen/logrus"
)

// WatchCommand for checking for new posts periodically and post them to twitter
type WatchCommand struct {
	Period int `short:"p" long:"period" default:"43200" description:"Period of checking for new posts in seconds"`
	CommonOptions
}

// SetCommonOptions sets common options in command
func (cmd *WatchCommand) SetCommonOptions(opts *CommonOptions) {
	cmd.CommonOptions = *opts
}

// Execute command method for flags.Commander
func (cmd *WatchCommand) Execute(args []string) error {
	log.Infof("%d\n", cmd.Period)
	return nil
}

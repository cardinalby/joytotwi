package main

import (
	"fmt"
	"github.com/imdario/mergo"
	"github.com/jessevdk/go-flags"
	"joytotwi/app/cmd"
	"joytotwi/app/cmd/common"
)

// merge common options from env with:
// 		* json config if appOpts.ConfigFile is set
// 		* appOpts itself otherwise
// and write them to command
func readCommonOptions(appOpts *cmd.AppOptions, command flags.Commander) error {
	opts, err := common.CreateOptionsFromEnv()
	if err != nil {
		return fmt.Errorf("error reading options from env variables. %s", err.Error())
	}

	var runOpts *common.Options
	if appOpts.ConfigFile != "" {
		runOpts, err = common.CreateOptionsFromJSONFile(appOpts.ConfigFile)
		if err != nil {
			return fmt.Errorf("error reading options from json config. %s", err.Error())
		}
	} else {
		runOpts = appOpts.GetCommonOptions()
	}

	err = mergo.Merge(opts, runOpts, mergo.WithOverride)
	if err != nil {
		return fmt.Errorf("error merging options. %s", err.Error())
	}

	err = opts.Validate()
	if err != nil {
		return fmt.Errorf("options validating error. %s", err.Error())
	}

	commonOptCommander, _ := command.(common.Commander)
	commonOptCommander.SetCommonOptions(opts)
	return nil
}

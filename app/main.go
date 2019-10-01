package main

import (
	"fmt"
	"joytotwi/app/cmd"
	"joytotwi/app/cmd/common"
	"os"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

func main() {
	var opts cmd.AppOptions
	envCommonOptions := common.Options{}
	envCommonOptErr := envCommonOptions.ReadFromEnv()

	p := flags.NewParser(&opts, flags.Default)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		if envCommonOptErr != nil {
			return fmt.Errorf("Error reading options from env variables: %s", envCommonOptErr.Error())
		}
		// read from json config if opts.ConfigFile is set or from opts and pass to command
		err := processCommonOptions(envCommonOptions, &opts, command)
		if err != nil {
			return err
		}

		err = command.Execute(args)
		if err != nil {
			log.Error(err)
			return fmt.Errorf("Command '%s' finished with error", p.Active.Name)
		}
		return err
	}
	parseFlags(p)
}

func processCommonOptions(commonOptions common.Options, appOpts *cmd.AppOptions, command flags.Commander) error {
	if appOpts.ConfigFile != "" {
		jsonErr := commonOptions.ReadFromJSONFile(appOpts.ConfigFile)
		if jsonErr != nil {
			return jsonErr
		}
	} else {
		commonOptions = appOpts.GetCommonOptions()
	}
	validateErr := commonOptions.Validate()
	if validateErr != nil {
		return validateErr
	}

	commonOptCommander, _ := command.(common.Commander)
	commonOptCommander.SetCommonOptions(&commonOptions)
	return nil
}

func parseFlags(p *flags.Parser) {
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
			}
			log.Fatal(err)
		} else {
			os.Exit(1)
		}
	}
}

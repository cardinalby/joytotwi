package main

import (
	"fmt"
	"joytotwi/app/cmd"
	"joytotwi/app/utils/stdemuxerhook"
	"os"

	"github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.AddHook(stdemuxerhook.New(log.StandardLogger()))

	var opts cmd.AppOptions

	p := flags.NewParser(&opts, flags.HelpFlag|flags.PassDoubleDash)
	p.CommandHandler = func(command flags.Commander, args []string) error {
		err := readCommonOptions(&opts, command)
		if err != nil {
			log.Fatal(err)
			return err
		}

		err = command.Execute(args)
		if err != nil {
			log.Error(err)
			log.Fatalf("command '%s' finished with error", p.Active.Name)
		}
		return nil
	}
	parseFlags(p)
}

func parseFlags(p *flags.Parser) {
	if _, err := p.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok {
			if flagsErr.Type == flags.ErrHelp {
				fmt.Println(flagsErr.Message)
				os.Exit(0)
			} else {
				log.Fatal(err)
				os.Exit(1)
			}
		}
	}
}

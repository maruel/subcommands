// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// sample-complex - Sample app to demonstrate example usage of package
// subcommand.
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/maruel/subcommands"
)

var application = &subcommands.DefaultApplication{
	Name:  "sample-complex",
	Title: "Sample tool to act as a skeleton for subcommands usage.",
	// Commands will be shown in this exact order, so you'll likely want to put
	// them in alphabetical order or in logical grouping.
	Commands: []*subcommands.Command{
		subcommands.Section("Nonsleepy commands."),
		cmdGreet,
		cmdHelp,
		cmdAsk,
		subcommands.Section("Sleepy commands."),
		cmdSleep,
	},
	EnvVars: map[string]subcommands.EnvVarDefinition{
		"GREET_STYLE": {
			ShortDesc: "Controls the type of greeting.",
			Default:   "Hi",
		},
		"VERBOSE_DREAMS": {
			Advanced:  true,
			ShortDesc: `If set to "1", shows dream while sleeping.`,
		},
	},
}

// cmdHelp overrides subcommands.CmdHelp to enable capture of the -advanced
// flag, which will be necessary when using subcommands.
var cmdHelp = &subcommands.Command{
	UsageLine: subcommands.CmdHelp.UsageLine,
	ShortDesc: subcommands.CmdHelp.ShortDesc,
	LongDesc:  subcommands.CmdHelp.LongDesc,
	CommandRun: func() subcommands.CommandRun {
		// Use the original implementation then steal -advanced.
		ret := subcommands.CmdHelp.CommandRun()
		ret.GetFlags().VisitAll(func(f *flag.Flag) {
			if f.Name == "advanced" {
				helpAdvanced = f.Value
			}
		})
		return ret
	},
}

// Warning: this is not concurrent safe. Only a concern when unit testing.
var helpAdvanced flag.Value

type sampleComplexApplication struct {
	*subcommands.DefaultApplication
	log *log.Logger
}

func main() {
	subcommands.KillStdLog()
	s := &sampleComplexApplication{application, log.New(ioutil.Discard, "", log.LstdFlags|log.Lmicroseconds)}
	os.Exit(subcommands.Run(s, nil))
}

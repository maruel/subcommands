// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/maruel/subcommands"
)

var askApplication = &subcommands.DefaultApplication{
	Name:  "sample-complex ask",
	Title: "Ask stuff.",
	// Commands will be shown in this exact order, so you'll likely want to put
	// them in alphabetical order or in logical grouping.
	Commands: []*subcommands.Command{
		cmdAskApple,
		cmdAskBeer,
		cmdAskArbitrary,
		cmdHelp,
	},
}

var cmdAsk = &subcommands.Command{
	UsageLine: "ask <subcommand>",
	ShortDesc: "asks questions",
	LongDesc:  "Asks one of the known subquestion.",
	CommandRun: func() subcommands.CommandRun {
		c := &askRun{}
		c.init()
		app := sampleComplexApplication{askApplication, nil}
		c.Flags.Usage = func() {
			advanced := helpAdvanced != nil && helpAdvanced.String() == "true"
			subcommands.Usage(os.Stderr, &app, advanced)
		}
		return c
	},
}

type askRun struct {
	commonFlags
}

type askCommonFlags struct {
	subcommands.CommandRunBase
}

func (a *askCommonFlags) init() {
}

func (a *askCommonFlags) parse(*sampleComplexApplication) error {
	return nil
}

func (c *askRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	d := a.(*sampleComplexApplication)
	// Create an inner application.
	app := sampleComplexApplication{askApplication, d.log}
	return subcommands.Run(&app, args)
}

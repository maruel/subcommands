// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"github.com/maruel/subcommands"
)

var cmdAsk = &subcommands.Command{
	UsageLine: "ask <subcommand>",
	ShortDesc: "asks questions",
	LongDesc:  "Asks one of the known subquestion.",
	CommandRun: func() subcommands.CommandRun {
		c := &askRun{}
		c.init()
		return c
	},
}

type askRun struct {
	commonFlags
}

// 'ask' is itself an application with subcommands.
type askApplication struct {
	sampleApplication
}

func (q askApplication) GetName() string {
	return q.sampleApplication.GetName() + " ask"
}

func (q askApplication) GetCommands() []*subcommands.Command {
	// Keep in alphabetical order of their name.
	return []*subcommands.Command{
		cmdAskApple,
		cmdAskBeer,
		subcommands.CmdHelp,
	}
}

func (c *askRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	d := a.(sampleApplication)
	// Create an inner application.
	return subcommands.Run(askApplication{d}, args)
}

// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/maruel/subcommands"
)

var cmdAsk = &subcommands.Command{
	UsageLine: "ask <subcommand>",
	ShortDesc: "asks questions",
	LongDesc:  "Asks one of the known subquestion.",
	CommandRun: func() subcommands.CommandRun {
		c := &askRun{}
		c.Init()
		return c
	},
}

type askRun struct {
	CommonFlags
	who string
}

func (c *askRun) main(a SampleApplication, args []string) error {
	if err := c.Parse(a, false); err != nil {
		return err
	}
	fmt.Fprintf(a.GetOut(), "TODO: Implement me!\n")
	return nil
}

// 'ask' is itself an application with subcommands.
type askApplication struct {
	SampleApplication
}

func (q askApplication) GetName() string {
	return q.SampleApplication.GetName() + " ask"
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
	d := a.(SampleApplication)
	// Create an inner application.
	return subcommands.Run(askApplication{d}, args)
}

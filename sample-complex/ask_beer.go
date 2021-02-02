// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"

	"github.com/maruel/subcommands"
)

var cmdAskBeer = &subcommands.Command{
	UsageLine: "beer <options>",
	ShortDesc: "asks for beer",
	LongDesc:  "Asks for beer.",
	CommandRun: func() subcommands.CommandRun {
		c := &askBeerRun{}
		c.init()
		c.Flags.StringVar(&c.file, "file", "", "Sets a new version of start_slave.py")
		return c
	},
}

type askBeerRun struct {
	commonFlags
	file string
}

func (c *askBeerRun) main(a askApplication) error {
	if err := c.parse(a, false); err != nil {
		return err
	}
	// This makes the process returns 1.
	return errors.New("it's a BYOB part")
}

func (c *askBeerRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if len(args) != 0 {
		fmt.Fprintf(a.GetErr(), "%s: Unknown arguments.\n", a.GetName())
		return 1
	}
	d := a.(askApplication)
	if err := c.main(d); err != nil {
		fmt.Fprintf(a.GetErr(), "%s: %s\n", a.GetName(), err)
		return 1
	}
	return 0
}

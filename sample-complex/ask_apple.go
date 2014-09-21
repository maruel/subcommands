// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/maruel/subcommands"
)

var cmdAskApple = &subcommands.Command{
	UsageLine: "apple <options>",
	ShortDesc: "asks for an apple",
	LongDesc:  "Asks for an apple.",
	CommandRun: func() subcommands.CommandRun {
		c := &askAppleRun{}
		c.Init()
		c.Flags.BoolVar(&c.bare, "bare", false, "Shows only the bot id, no meta data")
		return c
	},
}

type askAppleRun struct {
	CommonFlags
	bare bool
}

func (c *askAppleRun) main(a askApplication) error {
	// This command ignores -verbose.
	if err := c.Parse(a, true); err != nil {
		return err
	}
	fmt.Fprintf(a.GetOut(), "TODO: Implement 'ask apple'!\n")
	return nil
}

func (c *askAppleRun) Run(a subcommands.Application, args []string) int {
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

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
		c.init()
		c.Flags.BoolVar(&c.direct, "direct", false, "Be more direct")
		return c
	},
}

type askAppleRun struct {
	askCommonFlags
	direct bool
}

func (c *askAppleRun) main(a *sampleComplexApplication) error {
	if err := c.parse(a); err != nil {
		return err
	}
	if c.direct {
		fmt.Fprintf(a.GetOut(), "No way!\n")
		return nil
	}
	fmt.Fprintf(a.GetOut(), "Maybe one day.\n")
	return nil
}

func (c *askAppleRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if len(args) != 0 {
		fmt.Fprintf(a.GetErr(), "%s: Unknown arguments.\n", a.GetName())
		return 1
	}
	d := a.(*sampleComplexApplication)
	if err := c.main(d); err != nil {
		fmt.Fprintf(a.GetErr(), "%s: %s\n", a.GetName(), err)
		return 1
	}
	return 0
}

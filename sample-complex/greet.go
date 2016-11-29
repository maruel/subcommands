// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"

	"github.com/maruel/subcommands"
)

var cmdGreet = &subcommands.Command{
	UsageLine: "greet <who>",
	ShortDesc: "greets someone",
	LongDesc:  "Greets someone. This command has no specific option except the common ones.",
	CommandRun: func() subcommands.CommandRun {
		c := &greetRun{}
		c.Init()
		return c
	},
}

type greetRun struct {
	CommonFlags
}

func (c *greetRun) main(a SampleApplication, who, greeting string) error {
	if err := c.Parse(a, false); err != nil {
		return err
	}
	a.GetLog().Printf("Unnecessary logging, use -verbose to see it")
	fmt.Fprintf(a.GetOut(), "%s %s!\n", greeting, who)
	return nil
}

func (c *greetRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if len(args) != 1 {
		fmt.Fprintf(a.GetErr(), "%s: Can only greet one person at a time.\n", a.GetName())
		return 1
	}
	d := a.(SampleApplication)
	if err := c.main(d, args[0], env["GREET_STYLE"].Value); err != nil {
		fmt.Fprintf(a.GetErr(), "%s: %s\n", a.GetName(), err)
		return 1
	}
	return 0
}

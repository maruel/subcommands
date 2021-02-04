// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"

	"github.com/maruel/subcommands"
)

var cmdGreet = &subcommands.Command{
	UsageLine: "greet <who>",
	ShortDesc: "greets someone",
	LongDesc:  "Greets someone. This command has no specific option except the common ones.",
	CommandRun: func() subcommands.CommandRun {
		c := &greetRun{}
		c.init()
		return c
	},
}

type greetRun struct {
	commonFlags
}

func (c *greetRun) main(a *sampleComplexApplication, who, greeting string) error {
	if err := c.parse(a); err != nil {
		return err
	}
	log.Printf("Unnecessary logging, use -verbose to see it")
	fmt.Fprintf(a.GetOut(), "%s %s!\n", greeting, who)
	return nil
}

func (c *greetRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if len(args) != 1 {
		fmt.Fprintf(a.GetErr(), "%s: Can only greet one person at a time.\n", a.GetName())
		return 1
	}
	d := a.(*sampleComplexApplication)
	if err := c.main(d, args[0], env["GREET_STYLE"].Value); err != nil {
		fmt.Fprintf(a.GetErr(), "%s: %s\n", a.GetName(), err)
		return 1
	}
	return 0
}

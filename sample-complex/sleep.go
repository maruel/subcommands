// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/maruel/subcommands"
)

var cmdSleep = &subcommands.Command{
	UsageLine: "sleep <options>",
	ShortDesc: "sleeps for some time",
	LongDesc:  "Sleeps for some time, as desired.",
	CommandRun: func() subcommands.CommandRun {
		c := &sleepRun{}
		c.Flags.IntVar(&c.duration, "duration", 0, "Duration in seconds")
		return c
	},
}

type sleepRun struct {
	// This command doesn't implement the common flags.
	subcommands.CommandRunBase
	duration int
}

func (c *sleepRun) main(a sampleApplication, dream bool) error {
	if c.duration <= 0 {
		return errors.New("-duration is required")
	}
	fmt.Fprintf(a.GetOut(), "Sleeping for %ds.\n", c.duration)
	duration := time.Duration(c.duration) * time.Second
	if dream {
		chunk := time.Millisecond * 100
		for duration > 0 {
			fmt.Println("dreaming of sheep")
			time.Sleep(chunk)
			duration -= chunk
		}
	} else {
		time.Sleep(duration)
	}
	return nil
}

func (c *sleepRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if len(args) != 0 {
		fmt.Fprintf(a.GetErr(), "%s: Unsupported arguments.\n", a.GetName())
		return 1
	}
	d := a.(sampleApplication)
	// This main() wrapping simplifies the surfacing of errors into printing to
	// stderr then exiting with 1.
	if err := c.main(d, env["VERBOSE_DREAMS"].Value == "1"); err != nil {
		fmt.Fprintf(a.GetErr(), "%s: %s\n", a.GetName(), err)
		return 1
	}
	return 0
}

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
		c.Flags.DurationVar(&c.duration, "duration", time.Second, "Duration")
		return c
	},
}

type sleepRun struct {
	// This command doesn't implement the common flags.
	subcommands.CommandRunBase
	duration time.Duration
}

func (c *sleepRun) main(a *sampleComplexApplication, dream bool) error {
	if c.duration <= 0 {
		return errors.New("-duration is required")
	}
	fmt.Fprintf(a.GetOut(), "Sleeping for %v.\n", c.duration)
	duration := c.duration
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
	d := a.(*sampleComplexApplication)
	// This main() wrapping simplifies the surfacing of errors into printing to
	// stderr then exiting with 1.
	if err := c.main(d, env["VERBOSE_DREAMS"].Value == "1"); err != nil {
		fmt.Fprintf(a.GetErr(), "%s: %s\n", a.GetName(), err)
		return 1
	}
	return 0
}

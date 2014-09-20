// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// sample-simple - Sample app to demonstrate a very basic example usage of
// package subcommand.
//
// It implements 2 commands, one using an argument, the other using a flag.
// Help pages are automatically generated.
// The test cases cannot use t.Parallel() due to the use of global variables;
// log.*, os.Stdout and os.Stderr.
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/maruel/subcommands"
)

var application = &subcommands.DefaultApplication{
	Name:  "sample-simple",
	Title: "Sample tool to act as a skeleton for subcommands usage.",
	// Commands will be shown in this exact order, so you'll likely want to put
	// them in alphabetical order or in logical grouping.
	Commands: []*subcommands.Command{
		cmdGreet,
		cmdSleep,
		subcommands.CmdHelp,
	},
}

var cmdGreet = &subcommands.Command{
	UsageLine: "greet <who>",
	ShortDesc: "greets someone",
	LongDesc:  "Greets someone. This command has no specific option except the common ones.",
	CommandRun: func() subcommands.CommandRun {
		return &greetRun{}
	},
}

type greetRun struct {
	subcommands.CommandRunBase
}

func (c *greetRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "%s: Can only greet one person at a time.\n", a.GetName())
		return 1
	}
	fmt.Printf("Hi %s!\n", args[0])
	return 0
}

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
	subcommands.CommandRunBase
	duration int
}

func (c *sleepRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "%s: Unsupported arguments.\n", a.GetName())
		return 1
	}
	if c.duration <= 0 {
		fmt.Fprintf(os.Stderr, "%s: -duration is required.\n", a.GetName())
		return 1
	}
	log.Printf("Simulating sleeping for %ds.\n", c.duration)
	return 0
}

func main() {
	// It is not used Application.Logger.
	log.SetFlags(log.Lmicroseconds)
	os.Exit(subcommands.Run(application, nil))
}

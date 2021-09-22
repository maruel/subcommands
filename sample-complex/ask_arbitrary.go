// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/maruel/subcommands"
)

var cmdAskArbitrary = &subcommands.Command{
	UsageLine: "arbitrary <anything>",
	ShortDesc: "asks for anything you want",
	LongDesc:  "Asks for arbitrary arguments.",
	CommandRun: func() subcommands.CommandRun {
		// note that askArbitraryRun has no Flags
		return &askArbitraryRun{}
	},
}

type askArbitraryRun struct {
}

func (c *askArbitraryRun) GetFlags() *flag.FlagSet { return nil }

func (c *askArbitraryRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
	if len(args) == 0 {
		fmt.Fprintf(a.GetErr(), "%s: expected a question.", a.GetName())
		return 1
	}
	if last := args[len(args)-1]; !strings.HasSuffix(last, "?") {
		fmt.Fprintf(a.GetErr(), "%s: expected a question ending with `?`.", a.GetName())
		return 1
	}

	fmt.Println("You asked:", strings.Join(args, " "))
	fmt.Println("That's a great question!")
	return 0
}

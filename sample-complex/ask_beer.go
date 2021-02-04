// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/maruel/subcommands"
)

var cmdAskBeer = &subcommands.Command{
	UsageLine: "beer <options>",
	ShortDesc: "asks for beer",
	LongDesc:  "Asks for beer.",
	Advanced:  true,
	CommandRun: func() subcommands.CommandRun {
		c := &askBeerRun{}
		c.init()
		c.Flags.StringVar(&c.brand, "brand", "", "Which brand do you want?")
		return c
	},
}

type askBeerRun struct {
	askCommonFlags
	brand string
}

func (c *askBeerRun) main(a *sampleComplexApplication) error {
	if err := c.parse(a); err != nil {
		return err
	}
	if c.brand != "" && strings.ToLower(c.brand) != "unibroue" {
		fmt.Fprintf(a.GetOut(), "%q sounds interesting but we are partial to Unibroue.\n", c.brand)
		return nil
	}
	return errors.New("it's a BYOB part")
}

func (c *askBeerRun) Run(a subcommands.Application, args []string, env subcommands.Env) int {
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

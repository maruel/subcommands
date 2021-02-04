// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/maruel/subcommands"
)

type commonFlags struct {
	subcommands.CommandRunBase
	Verbose bool
}

func (c *commonFlags) init() {
	c.Flags.BoolVar(&c.Verbose, "verbose", false, "Enable verbose output.")
}

func (c *commonFlags) parse(a *sampleComplexApplication) error {
	if c.Verbose {
		// Enable logging when -verbose is specified.
		a.log = log.New(a.GetErr(), a.log.Prefix(), a.log.Flags())
	}
	return nil
}

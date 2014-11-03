// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"log"

	"github.com/maruel/subcommands"
)

// Common flags.
type CommonFlags struct {
	subcommands.CommandRunBase
	Verbose bool
}

func (c *CommonFlags) Init() {
	c.Flags.BoolVar(&c.Verbose, "verbose", false, "Enable verbose output.")
}

func (c *CommonFlags) Parse(d SampleApplication, special bool) error {
	if c.Verbose && !special {
		// Enable logging when -verbose is specified.
		a := d.(*sample)
		a.log = log.New(d.GetErr(), a.log.Prefix(), a.log.Flags())
	}
	return nil
}

// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// sample - Sample app to demonstrate example usage of package subcommand.
package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/maruel/subcommands"
	"github.com/maruel/subcommands/subcommandstest"
)

var application = &subcommands.DefaultApplication{
	Name:  "sample",
	Title: "Sample tool to act as a skeleton for subcommands usage.",
	// Commands will be shown in this exact order, so you'll likely want to put
	// them in alphabetical order or in logical grouping.
	Commands: []*subcommands.Command{
		cmdGreet,
		subcommands.CmdHelp,
		cmdAsk,
		cmdSleep,
	},
}

type SampleApplication interface {
	// TODO(maruel): This is wrong, subcommandtest should only be referenced in
	// unit tests. Figure out a way to better plug logging.
	subcommandstest.Application

	// Add anything desired, in particular if you'd like to crete a fake
	// application during testing.
}

type sample struct {
	*subcommands.DefaultApplication
	log *log.Logger
}

// GetLog implements subcommandstest.Application.
func (s *sample) GetLog() *log.Logger {
	return s.log
}

func main() {
	subcommands.KillStdLog()
	s := &sample{application, log.New(ioutil.Discard, "", log.LstdFlags|log.Lmicroseconds)}
	os.Exit(subcommands.Run(s, nil))
}

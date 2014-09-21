// Copyright 2012 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package subcommands

import (
	"log"
)

type panicWriter struct {
}

func (p panicWriter) Write(b []byte) (n int, err error) {
	panic("unexpected write")
}

// KillStdLog sets an output that will panic if used. This permits trapping any
// log.*() calls that would inhibit efficient use of t.Parallel().
func KillStdLog() {
	log.SetOutput(panicWriter{})
}

type nullWriter struct {
}

func (p nullWriter) Write(b []byte) (n int, err error) {
	return len(b), nil
}

// VoidStdLog sets an output that will be ignored. This permits ignoring any
// log.*() calls that would inhibit efficient use of t.Parallel().
func VoidStdLog() {
	log.SetOutput(panicWriter{})
}

// Copyright 2012 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package subcommands

import (
	"log"
)

// PanicWriter is an io.Writer that will panic if used.
type PanicWriter struct {
}

// Write implements io.Writer.
func (p PanicWriter) Write(b []byte) (n int, err error) {
	panic("unexpected write")
}

// KillStdLog sets an output that will panic if used. This permits trapping any
// log.*() calls that would inhibit efficient use of t.Parallel().
func KillStdLog() {
	log.SetOutput(PanicWriter{})
}

// NullWriter is an io.Writer that ignores everything written to it.
type NullWriter struct {
}

// Write implements io.Writer.
func (p NullWriter) Write(b []byte) (n int, err error) {
	return len(b), nil
}

// VoidStdLog sets an output that will be ignored. This permits ignoring any
// log.*() calls that would inhibit efficient use of t.Parallel().
func VoidStdLog() {
	log.SetOutput(NullWriter{})
}

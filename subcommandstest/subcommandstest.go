/* Copyright 2012 Marc-Antoine Ruel. Licensed under the Apache License, Version
2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law or
agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied. See the License for the specific language governing permissions and
limitations under the License. */

// Package subcommandstest includes tools to help with concurrent testing.
package subcommandstest

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/maruel/subcommands"
	"github.com/maruel/ut"
)

// Logging is a global object so it can't be checked for when tests are run in
// parallel.
var bufLog bytes.Buffer

// DisableLogOutput redirect log's default logger to a buffer and returns it.
// This function should be called in an init function.
func DisableLogOutput() *bytes.Buffer {
	log.SetOutput(&bufLog)
	log.SetFlags(log.Lmicroseconds)
	return &bufLog
}

// TB wraps a testing.T object and adds functionality specific to
// command_support.
//
// BUG: testing.TB is now a thing, so this struct should likely be renamed to
// reduce confusion?
type TB struct {
	*testing.T
	bufLog bytes.Buffer
	bufOut bytes.Buffer
	bufErr bytes.Buffer
	log    *log.Logger
}

// MakeTB returns a fully initialized TB instance.
func MakeTB(t *testing.T) *TB {
	tb := &TB{T: t}
	tb.log = log.New(&tb.bufLog, "", log.Lmicroseconds)
	return tb
}

// PrintIf trims the whitespace around a string encoded as []byte and prints it
// only if not empty.
func PrintIf(b []byte, name string) {
	s := strings.TrimSpace(string(b))
	if len(s) != 0 {
		fmt.Fprintf(os.Stderr, "\n\\/ \\/ %s \\/ \\/\n%s\n/\\ /\\ %s /\\ /\\\n", name, s, name)
	}
}

// CheckBuffer asserts the content of os.Stdout and os.Stderr mocks.
func (t *TB) CheckBuffer(out, err bool) {
	if out {
		// Print Stderr to see what happened.
		ut.AssertEqualf(t, true, t.bufOut.Len() != 0, "Expected stdout")
	} else {
		ut.AssertEqualf(t, t.bufOut.Len(), 0, "Unexpected stdout")
	}

	if err {
		ut.AssertEqualf(t, true, t.bufErr.Len() != 0, "Expected stderr")
	} else {
		ut.AssertEqualf(t, t.bufErr.Len(), 0, "Unexpected stderr")
	}
	t.bufOut.Reset()
	t.bufErr.Reset()
}

// CheckOut asserts that what was printed out Application.GetOut() matches what
// is expected.
// TODO(maruel): It doesn't matches the use case where the match must be fuzzy,
// for example when non-deterministic data is included in the output.
func (t *TB) CheckOut(expected string) {
	actual := t.bufOut.String()
	ut.AssertEqual(t, expected, actual)
	t.bufOut.Reset()
}

// Verbose sets the current context as verbose. It immediately prints out all
// logs generated for this specific test case up to now and redirects the log
// to os.Stderr so the following log is directly output to the console.
func (t *TB) Verbose() {
	if t.bufLog.Len() != 0 {
		_, _ = os.Stderr.Write(t.bufLog.Bytes())
	}
	t.log = log.New(os.Stderr, "", log.Lmicroseconds)
}

// GetLog implements Application.
func (t *TB) GetLog() *log.Logger {
	return t.log
}

// Application supports all of subcommands.Application and adds GetLog() for
// testing purposes.
type Application interface {
	subcommands.Application
	GetLog() *log.Logger
}

// ApplicationMock wrap both an Application and a TB. ApplicationMock
// implements GetOut and GetErr and adds GetLog(). GetLog() is implemented by
// TB.
type ApplicationMock struct {
	subcommands.Application
	*TB
}

// GetOut implements subcommands.Application.
func (a *ApplicationMock) GetOut() io.Writer {
	return &a.bufOut
}

// GetErr implements subcommands.Application.
func (a *ApplicationMock) GetErr() io.Writer {
	return &a.bufErr
}

// MakeAppMock returns an initialized ApplicationMock.
func MakeAppMock(t *testing.T, a subcommands.Application) *ApplicationMock {
	return &ApplicationMock{a, MakeTB(t)}
}

/* Copyright 2012 Marc-Antoine Ruel. Licensed under the Apache License, Version
2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law or
agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied. See the License for the specific language governing permissions and
limitations under the License. */

// Includes tools to help with concurrent testing.
package subcommandstest

import (
	"bytes"
	"fmt"
	"github.com/maruel/subcommands"
	"io"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"testing"
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

// ReduceStackTrace reduces the amount of data in a stack trace. It trims the
// first 2 lines and remove the file paths and function pointers to only keep
// the file names and line numbers.
func ReduceStackTrace(b []byte) []byte {
	lines := strings.Split(string(b), "\n")
	if len(lines) > 2 {
		lines = lines[2:]
	}
	for i := 0; i < len(lines); i++ {
		if !strings.HasPrefix(lines[i], "\t") {
			// /path/to/file.go:<lineno> (<addr>)
			// debug.Stack() uses "/" even on Windows.
			start := strings.LastIndex(lines[i], "/")
			end := strings.LastIndex(lines[i], " ")
			if start != -1 && end != -1 {
				lines[i] = lines[i][start+1 : end]
			}
		}
	}
	return []byte(strings.Join(lines, "\n"))
}

// Assertf prints the stack trace to ease debugging.  It's slightly slower than
// an explicit condition in the test but its more compact.
func (t *TB) Assertf(truth bool, format string, values ...interface{}) {
	if !truth {
		PrintIf(t.bufOut.Bytes(), "STDOUT")
		PrintIf(t.bufErr.Bytes(), "STDERR")
		PrintIf(t.bufLog.Bytes(), "LOG")
		os.Stderr.Write([]byte("\n"))
		os.Stderr.Write(ReduceStackTrace(debug.Stack()))
		t.Fatalf(format, values...)
	}
}

// CheckBuffer asserts the content of os.Stdout and os.Stderr mocks.
func (t *TB) CheckBuffer(out, err bool) {
	if out {
		// Print Stderr to see what happened.
		t.Assertf(t.bufOut.Len() != 0, "Expected stdout")
	} else {
		t.Assertf(t.bufOut.Len() == 0, "Unexpected stdout")
	}

	if err {
		t.Assertf(t.bufErr.Len() != 0, "Expected stderr")
	} else {
		t.Assertf(t.bufErr.Len() == 0, "Unexpected stderr")
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
	t.Assertf(expected == actual, "Expected:\n%s\nActual:\n%s", expected, actual)
	t.bufOut.Reset()
}

// Verbose sets the current context as verbose. It immediately prints out all
// logs generated for this specific test case up to now and redirects the log
// to os.Stderr so the following log is directly output to the console.
func (t *TB) Verbose() {
	if t.bufLog.Len() != 0 {
		os.Stderr.Write(t.bufLog.Bytes())
	}
	t.log = log.New(os.Stderr, "", log.Lmicroseconds)
}

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

func (a *ApplicationMock) GetOut() io.Writer {
	return &a.bufOut
}

func (a *ApplicationMock) GetErr() io.Writer {
	return &a.bufErr
}

// MakeAppMock returns an initialized ApplicationMock.
func MakeAppMock(t *testing.T, a subcommands.Application) *ApplicationMock {
	return &ApplicationMock{a, MakeTB(t)}
}

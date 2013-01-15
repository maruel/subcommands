/* Copyright 2012 Marc-Antoine Ruel. Licensed under the Apache License, Version
2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law or
agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied. See the License for the specific language governing permissions and
limitations under the License. */

package subcommandstest

import (
	"github.com/maruel/subcommands"
	"testing"
)

// Testing subcommands would require importing subcommandstest. To not create
// an import cycle between subcommands and subcommandstest, the public part of
// module subcommands is tested here.

func init() {
	DisableLogOutput()
}

func TestHelp(t *testing.T) {
	t.Parallel()
	app := &subcommands.DefaultApplication{
		Name:  "name",
		Title: "doc",
		Commands: []*subcommands.Command{
			subcommands.CmdHelp,
		},
	}
	a := MakeAppMock(t, app)
	args := []string{"help"}
	r := subcommands.Run(a, args)
	a.Assertf(r == 0, "Unexpected return code %d", r)
	a.CheckBuffer(true, false)
}

func TestHelpBadFlag(t *testing.T) {
	t.Parallel()
	app := &subcommands.DefaultApplication{
		Name:  "name",
		Title: "doc",
		Commands: []*subcommands.Command{
			subcommands.CmdHelp,
		},
	}
	a := MakeAppMock(t, app)
	args := []string{"help", "-foo"}
	r := subcommands.Run(a, args)
	a.Assertf(r == 2, "Unexpected return code %d", r)
	a.CheckBuffer(false, true)
}

func TestHelpBadCommand(t *testing.T) {
	t.Parallel()
	app := &subcommands.DefaultApplication{
		Name:  "name",
		Title: "doc",
		Commands: []*subcommands.Command{
			subcommands.CmdHelp,
		},
	}
	a := MakeAppMock(t, app)
	args := []string{"help", "non_existing_command"}
	r := subcommands.Run(a, args)
	a.Assertf(r == 2, "Unexpected return code %d", r)
	a.CheckBuffer(false, true)
}

func TestBadCommand(t *testing.T) {
	t.Parallel()
	app := &subcommands.DefaultApplication{
		Name:  "name",
		Title: "doc",
		Commands: []*subcommands.Command{
			subcommands.CmdHelp,
		},
	}
	a := MakeAppMock(t, app)
	args := []string{"non_existing_command"}
	r := subcommands.Run(a, args)
	a.Assertf(r == 2, "Unexpected return code %d", r)
	a.CheckBuffer(false, true)
}

func TestReduceStackTrace(t *testing.T) {
	t.Parallel()
	tb := MakeTB(t)
	data := "/home/joe/gocode/src/github.com/maruel/dumbcas/command_support_test.go:93 (0x43acb9)\n" +
		"\tcom/maruel/dumbcas.(*TB).Assertf: os.Stderr.Write(ReduceStackTrace(debug.Stack()))\n" +
		"/home/joe/gocode/src/github.com/maruel/dumbcas/command_support_test.go:109 (0x43aeaf)\n" +
		"\tcom/maruel/dumbcas.(*TB).CheckBuffer: t.Assertf(t.bufErr.Len() != 0, \"Unexpected stderr\")\n" +
		"/home/joe/gocode/src/github.com/maruel/dumbcas/web_test.go:57 (0x440109)\n" +
		"\tcom/maruel/dumbcas.(*WebDumbcasAppMock).closeWeb: f.CheckBuffer(false, false)\n" +
		"/home/joe/gocode/src/github.com/maruel/dumbcas/web_test.go:147 (0x441a54)\n" +
		"\tcom/maruel/dumbcas.TestWeb: f.closeWeb()\n"

	// Much nicer!
	expected := "command_support_test.go:109\n" +
		"\tcom/maruel/dumbcas.(*TB).CheckBuffer: t.Assertf(t.bufErr.Len() != 0, \"Unexpected stderr\")\n" +
		"web_test.go:57\n" +
		"\tcom/maruel/dumbcas.(*WebDumbcasAppMock).closeWeb: f.CheckBuffer(false, false)\n" +
		"web_test.go:147\n" +
		"\tcom/maruel/dumbcas.TestWeb: f.closeWeb()\n"

	actual := string(ReduceStackTrace([]byte(data)))
	tb.Assertf(expected == actual, "ReduceStackTrace() failed parsing.\nActual:\n%s\n\nExpected:\n%s", expected, actual)
}

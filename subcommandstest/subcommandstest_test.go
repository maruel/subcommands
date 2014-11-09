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
	"testing"

	"github.com/maruel/subcommands"
	"github.com/maruel/ut"
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
	ut.AssertEqual(t, r, 0)
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
	ut.AssertEqual(t, r, 2)
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
	ut.AssertEqual(t, r, 2)
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
	ut.AssertEqual(t, r, 2)
	a.CheckBuffer(false, true)
}

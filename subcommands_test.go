// Copyright 2012 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package subcommands permits a Go application to implement subcommands support
// similar to what is supported by the 'go' tool.
//
// The library is designed so that the test cases can run concurrently.
// Using global flags variables is discouraged to keep your program testable
// conccurently.
package subcommands

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"testing"

	"github.com/maruel/ut"
)

func TestFindCommand(t *testing.T) {
	commands := []*Command{
		{UsageLine: "Fo"},
		{UsageLine: "Foo bar"},
		{UsageLine: "LongCommand"},
	}
	a := &DefaultApplication{Commands: commands}

	// Exact
	ut.AssertEqual(t, commands[0], FindCommand(a, "Fo"))
	ut.AssertEqual(t, commands[1], FindCommand(a, "Foo"))
	ut.AssertEqual(t, commands[2], FindCommand(a, "LongCommand"))

	// Prefix
	ut.AssertEqual(t, (*Command)(nil), FindCommand(a, "F"))
	ut.AssertEqual(t, (*Command)(nil), FindCommand(a, "LongC"))

	// Case insensitive
	ut.AssertEqual(t, (*Command)(nil), FindCommand(a, "fo"))
	ut.AssertEqual(t, (*Command)(nil), FindCommand(a, "foo"))
	ut.AssertEqual(t, (*Command)(nil), FindCommand(a, "longcommand"))
}

func TestFindNearestCommand(t *testing.T) {
	commands := []*Command{
		{UsageLine: "Fo"},
		{UsageLine: "Foo"},
		{UsageLine: "LongCommand"},
		{UsageLine: "LargCommand"},
		Section("bar"),
	}
	a := &DefaultApplication{Commands: commands}

	// Exact
	ut.AssertEqual(t, commands[0], FindNearestCommand(a, "Fo"))
	ut.AssertEqual(t, commands[1], FindNearestCommand(a, "Foo"))
	ut.AssertEqual(t, commands[2], FindNearestCommand(a, "LongCommand"))

	// Prefix
	ut.AssertEqual(t, (*Command)(nil), FindNearestCommand(a, "F"))
	ut.AssertEqual(t, commands[2], FindNearestCommand(a, "Lo"))

	// Case insensitive
	ut.AssertEqual(t, (*Command)(nil), FindNearestCommand(a, "fo"))
	ut.AssertEqual(t, commands[1], FindNearestCommand(a, "foo"))
	ut.AssertEqual(t, commands[2], FindNearestCommand(a, "longcommand"))
	ut.AssertEqual(t, commands[2], FindNearestCommand(a, "longc"))

	// Based on levenshtein distance
	ut.AssertEqual(t, (*Command)(nil), FindNearestCommand(a, "Fof"))
	ut.AssertEqual(t, commands[2], FindNearestCommand(a, "LongCommandd"))
	ut.AssertEqual(t, commands[2], FindNearestCommand(a, "LongCmomand"))
	ut.AssertEqual(t, commands[2], FindNearestCommand(a, "ongCommand"))
	ut.AssertEqual(t, (*Command)(nil), FindNearestCommand(a, "LangCommand"))

	// Section cannot be found.
	ut.AssertEqual(t, (*Command)(nil), FindNearestCommand(a, "bar"))
}

func TestHelpOutput(t *testing.T) {
	a := &DefaultApplication{
		Commands: []*Command{
			{UsageLine: "Foo", ShortDesc: "A foo"},
			{UsageLine: "SuperDuperLongLine", ShortDesc: "A long thing", Advanced: true},
		},
		EnvVars: map[string]EnvVarDefinition{
			"EVAR":            {ShortDesc: "Desc"},
			"SUPER_LONG_EVAR": {ShortDesc: "Desc", Advanced: true},
			"DFLT_EVAR":       {ShortDesc: "Desc", Default: "yep"},
		},
	}

	buf := bytes.Buffer{}
	usage(&buf, a, false)

	ut.AssertEqual(t, buf.String(), `

Usage:   [command] [arguments]

Commands:
  Foo  A foo

Environment Variables:
  DFLT_EVAR  Desc (Default: "yep")
  EVAR       Desc


Use " help [command]" for more information about a command.
Use " help -advanced" to display all commands.

`)

	buf.Reset()
	usage(&buf, a, true)
	ut.AssertEqual(t, buf.String(), `

Usage:   [command] [arguments]

Commands:
  Foo                 A foo
  SuperDuperLongLine  A long thing

Environment Variables:
  DFLT_EVAR        Desc (Default: "yep")
  EVAR             Desc
  SUPER_LONG_EVAR  Desc


Use " help [command]" for more information about a command.

`)

}

func TestDefaultApplication_GetOut_GetErr(t *testing.T) {
	a := DefaultApplication{}
	ut.AssertEqual(t, a.GetOut().(*os.File), os.Stdout)
	ut.AssertEqual(t, a.GetErr().(*os.File), os.Stderr)
}

func TestCommandRunBase_GetFlags(t *testing.T) {
	c := CommandRunBase{}
	ut.AssertEqual(t, c.GetFlags(), &c.Flags)
}

func TestCmdHelp(t *testing.T) {
	data := []struct {
		args []string
		out  string
		err  string
		exit int
	}{
		{
			[]string{"help"},
			"Title\n" +
				"\n" +
				"Usage:  App [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"  help  prints help about a command\n" +
				"\n" +
				"\n" +
				"Use \"App help [command]\" for more information about a command.\n" +
				"Use \"App help -advanced\" to display all commands.\n" +
				"\n",
			"",
			0,
		},
		{
			[]string{"help", "-advanced"},
			"Title\n" +
				"\n" +
				"Usage:  App [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"  help  prints help about a command\n" +
				"  foo   foo\n" +
				"\n" +
				"\n" +
				"Use \"App help [command]\" for more information about a command.\n" +
				"\n",
			"",
			0,
		},
		{
			[]string{"help", "foo"},
			"",
			"Foo.\n" +
				"\n" +
				"usage:  App foo\n",
			0,
		},
		{
			[]string{"help", "foo", "bar"},
			"",
			"App: Too many arguments given\n" +
				"\n" +
				"Run 'App help' for usage.\n",
			2,
		},
		{
			[]string{"foo", "-help"},
			"",
			"Foo.\n" +
				"\n" +
				"usage:  App foo\n",
			2,
		},
		{
			[]string{"help", "inexistant"},
			"",
			"App: unknown command `inexistant`\n" +
				"\n" +
				"Run 'App help' for usage.\n",
			2,
		},
		{
			[]string{"inexistant"},
			"",
			"App: unknown command `inexistant`\n" +
				"\n" +
				"Run 'App help' for usage.\n",
			2,
		},
		{
			nil,
			"",
			"Title\n" +
				"\n" +
				"Usage:  App [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"  help  prints help about a command\n" +
				"\n" +
				"\n" +
				"Use \"App help [command]\" for more information about a command.\n" +
				"Use \"App help -advanced\" to display all commands.\n" +
				"\n",
			2,
		},
	}

	for i, line := range data {
		line := line
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			a := application{
				DefaultApplication: DefaultApplication{
					Name:  "App",
					Title: "Title",
					Commands: []*Command{
						CmdHelp,
						{
							UsageLine: "foo",
							ShortDesc: "foo",
							LongDesc:  "Foo.",
							Advanced:  true,
							CommandRun: func() CommandRun {
								return &command{}
							},
						},
					},
				},
			}
			ut.AssertEqual(t, Run(&a, line.args), line.exit)
			ut.AssertEqual(t, a.out.String(), line.out)
			ut.AssertEqual(t, a.err.String(), line.err)
		})
	}
}

type application struct {
	DefaultApplication
	out bytes.Buffer
	err bytes.Buffer
}

func (a *application) GetOut() io.Writer {
	return &a.out
}

func (a *application) GetErr() io.Writer {
	return &a.err
}

type command struct {
	CommandRunBase
}

func (c *command) Run(a Application, args []string, env Env) int {
	return 42
}

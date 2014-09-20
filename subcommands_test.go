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
	"github.com/maruel/ut"
	"testing"
)

func TestFindCommand(t *testing.T) {
	commands := []*Command{
		&Command{UsageLine: "Fo"},
		&Command{UsageLine: "Foo bar"},
		&Command{UsageLine: "LongCommand"},
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
		&Command{UsageLine: "Fo"},
		&Command{UsageLine: "Foo"},
		&Command{UsageLine: "LongCommand"},
		&Command{UsageLine: "LargCommand"},
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
	ut.AssertEqual(t, commands[2], FindNearestCommand(a, "ongCommand"))
	ut.AssertEqual(t, (*Command)(nil), FindNearestCommand(a, "LangCommand"))
}

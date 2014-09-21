// Copyright 2014 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"github.com/maruel/subcommands"
	"github.com/maruel/ut"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// t.Parallel() cannot be used here, see main.go for rationale.
// In addition, logging must be zapped out.

// TODO(maruel): Create an in-memory os.File, couldn't quickly find a ready
// made fake only.
func newTempFile(t testing.TB) *os.File {
	f, err := ioutil.TempFile("", "sample-simple")
	ut.AssertEqual(t, nil, err)
	return f
}

// mockStdout mocks os.Stdout manually. To have it mocked automatically, see
// sample-complex.
func mockStdout(t testing.TB) func() {
	oldStdout := os.Stdout
	os.Stdout = newTempFile(t)
	return func() {
		os.Stdout.Close()
		os.Stdout = oldStdout
	}
}

func assertStdout(t testing.TB, expected string) {
	os.Stdout.Seek(0, 0)
	actual, err := ioutil.ReadAll(os.Stdout)
	ut.AssertEqual(t, nil, err)
	ut.AssertEqual(t, expected, string(actual))
}

func TestGreet(t *testing.T) {
	defer mockStdout(t)()

	ut.AssertEqual(t, 0, subcommands.Run(application, []string{"greet", "active tester"}))
	assertStdout(t, "Hi active tester!\n")
}

func TestSleep(t *testing.T) {
	defer mockStdout(t)()

	// If running with "go test -v", the following log entry will be printed:
	// utiltest.go:132: 2010/01/02 03:04:05 Simulating sleeping for 1s.
	out := ut.NewWriter(t)
	defer out.Close()
	log.SetOutput(out)
	ut.AssertEqual(t, 0, subcommands.Run(application, []string{"sleep", "-duration", "1"}))
	assertStdout(t, "")
}

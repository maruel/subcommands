// Copyright 2021 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// "go run ." on before 1.11 has a quirk where it runs but prints
// go run: no go files listed

//go:build go1.11
// +build go1.11

package main_test

import (
	"bytes"
	"os/exec"
	"strconv"
	"syscall"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHelp(t *testing.T) {
	data := []struct {
		args     []string
		expected string
		exitCode int
	}{
		{
			[]string{"-help"},
			"Sample tool to act as a skeleton for subcommands usage.\n" +
				"\n" +
				"Usage:  sample-complex [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"         \n" +
				"\tNonsleepy commands.\n" +
				"  greet  greets someone\n" +
				"  help   prints help about a command\n" +
				"  ask    asks questions\n" +
				"         \n" +
				"\tSleepy commands.\n" +
				"  sleep  sleeps for some time\n" +
				"\n" +
				"Environment Variables:\n" +
				"  GREET_STYLE  Controls the type of greeting. (Default: \"Hi\")\n" +
				"\n" +
				"\n" +
				"Use \"sample-complex help [command]\" for more information about a command.\n" +
				"Use \"sample-complex help -advanced\" to display all commands.\n" +
				"\n",
			0,
		},
		{
			[]string{"help"},
			"Sample tool to act as a skeleton for subcommands usage.\n" +
				"\n" +
				"Usage:  sample-complex [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"         \n" +
				"\tNonsleepy commands.\n" +
				"  greet  greets someone\n" +
				"  help   prints help about a command\n" +
				"  ask    asks questions\n" +
				"         \n" +
				"\tSleepy commands.\n" +
				"  sleep  sleeps for some time\n" +
				"\n" +
				"Environment Variables:\n" +
				"  GREET_STYLE  Controls the type of greeting. (Default: \"Hi\")\n" +
				"\n" +
				"\n" +
				"Use \"sample-complex help [command]\" for more information about a command.\n" +
				"Use \"sample-complex help -advanced\" to display all commands.\n" +
				"\n",
			0,
		},
		{
			[]string{"help", "-advanced"},
			"Sample tool to act as a skeleton for subcommands usage.\n" +
				"\n" +
				"Usage:  sample-complex [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"         \n" +
				"\tNonsleepy commands.\n" +
				"  greet  greets someone\n" +
				"  help   prints help about a command\n" +
				"  ask    asks questions\n" +
				"         \n" +
				"\tSleepy commands.\n" +
				"  sleep  sleeps for some time\n" +
				"\n" +
				"Environment Variables:\n" +
				"  GREET_STYLE     Controls the type of greeting. (Default: \"Hi\")\n" +
				"  VERBOSE_DREAMS  If set to \"1\", shows dream while sleeping." +
				"\n" +
				"\n" +
				"\n" +
				"Use \"sample-complex help [command]\" for more information about a command.\n" +
				"\n",
			0,
		},

		{
			[]string{"-help", "ask"},
			"Ask stuff.\n" +
				"\n" +
				"Usage:  sample-complex ask [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"  apple      asks for an apple\n" +
				"  arbitrary  asks for anything you want\n" +
				"  help       prints help about a command\n" +
				"\n" +
				"\n" +
				"Use \"sample-complex ask help [command]\" for more information about a command.\n" +
				"Use \"sample-complex ask help -advanced\" to display all commands.\n" +
				"\n",
			0,
		},
		// TODO(maruel): This feels incorrect.
		{
			[]string{"ask", "-help"},
			"Ask stuff.\n" +
				"\n" +
				"Usage:  sample-complex ask [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"  apple      asks for an apple\n" +
				"  arbitrary  asks for anything you want\n" +
				"  help       prints help about a command\n" +
				"\n" +
				"\n" +
				"Use \"sample-complex ask help [command]\" for more information about a command.\n" +
				"Use \"sample-complex ask help -advanced\" to display all commands.\n" +
				"\n" +
				"exit status 2\n",
			1,
		},
		{
			[]string{"help", "ask"},
			"Ask stuff.\n" +
				"\n" +
				"Usage:  sample-complex ask [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"  apple      asks for an apple\n" +
				"  arbitrary  asks for anything you want\n" +
				"  help       prints help about a command\n" +
				"\n" +
				"\n" +
				"Use \"sample-complex ask help [command]\" for more information about a command.\n" +
				"Use \"sample-complex ask help -advanced\" to display all commands.\n" +
				"\n",
			0,
		},
		{
			[]string{"help", "-advanced", "ask"},
			"Ask stuff.\n" +
				"\n" +
				"Usage:  sample-complex ask [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"  apple      asks for an apple\n" +
				"  beer       asks for beer\n" +
				"  arbitrary  asks for anything you want\n" +
				"  help       prints help about a command\n" +
				"\n" +
				"\n" +
				"Use \"sample-complex ask help [command]\" for more information about a command.\n" +
				"\n",
			0,
		},
		{
			[]string{"ask", "help"},
			"Ask stuff.\n" +
				"\n" +
				"Usage:  sample-complex ask [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"  apple      asks for an apple\n" +
				"  arbitrary  asks for anything you want\n" +
				"  help       prints help about a command\n" +
				"\n" +
				"\n" +
				"Use \"sample-complex ask help [command]\" for more information about a command.\n" +
				"Use \"sample-complex ask help -advanced\" to display all commands.\n" +
				"\n",
			0,
		},
		{
			[]string{"ask", "help", "-advanced"},
			"Ask stuff.\n" +
				"\n" +
				"Usage:  sample-complex ask [command] [arguments]\n" +
				"\n" +
				"Commands:\n" +
				"  apple      asks for an apple\n" +
				"  beer       asks for beer\n" +
				"  arbitrary  asks for anything you want\n" +
				"  help       prints help about a command\n" +
				"\n" +
				"\n" +
				"Use \"sample-complex ask help [command]\" for more information about a command.\n" +
				"\n",
			0,
		},
		{
			[]string{"ask", "arbitrary", "-flags", "-don't", "matter?"},
			"You asked: -flags -don't matter?\nThat's a great question!\n",
			0,
		},
	}
	for i, line := range data {
		line := line
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			cmd := exec.Command("go", append([]string{"run", ".", "--"}, line.args...)...)
			buf := bytes.Buffer{}
			cmd.Stdout = &buf
			cmd.Stderr = &buf
			exitCode := 0
			if err := cmd.Run(); err != nil {
				if exiterr, ok := err.(*exec.ExitError); ok {
					if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
						exitCode = status.ExitStatus()
					}
				}
			}
			if diff := cmp.Diff(line.expected, buf.String()); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
			if exitCode != line.exitCode {
				t.Fatal(exitCode)
			}
		})
	}
}

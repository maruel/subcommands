// Copyright 2012 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package subcommands

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Logger defines a reduced interface that covers most of what the standard
// package log exposes in a more coherent way.
//
// The author of this code thinks Fatal*(), Panic*() and *ln() are highly
// redundants so they were not implemented.
type Logger interface {
	Flags() int
	Prefix() string
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	SetFlags(flag int)
	SetOutput(w io.Writer)
	SetPrefix(prefix string)
}

type privateLoggerImpl struct {
	l *log.Logger
}

// NewLogger returns a Logger instead with the current settings and output set
// to os.Stderr.
func NewLogger() Logger {
	return &privateLoggerImpl{log.New(os.Stderr, log.Prefix(), log.Flags())}
}

// NewPrivateLogger converts a log.Logger into a Logger interface.
func NewPrivateLogger(l *log.Logger) Logger {
	return &privateLoggerImpl{l}
}

func (l *privateLoggerImpl) Flags() int {
	return l.l.Flags()
}

func (l *privateLoggerImpl) Prefix() string {
	return l.l.Prefix()
}

func (l *privateLoggerImpl) Print(v ...interface{}) {
	l.l.Output(2, fmt.Sprint(v...))
}

func (l *privateLoggerImpl) Printf(format string, v ...interface{}) {
	l.l.Output(2, fmt.Sprintf(format, v...))
}

func (l *privateLoggerImpl) SetFlags(flag int) {
	l.l.SetFlags(flag)
}

func (l *privateLoggerImpl) SetOutput(w io.Writer) {
	l.l = log.New(w, l.l.Prefix(), l.l.Flags())
}

func (l *privateLoggerImpl) SetPrefix(prefix string) {
	l.l.SetPrefix(prefix)
}

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

// VoidStdLog sets an output that will be ignored. This permits ignoring any
// log.*() calls that would inhibit efficient use of t.Parallel().
func VoidStdLog() {
	log.SetOutput(io.Discard)
}

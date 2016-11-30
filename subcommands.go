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
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/texttheater/golang-levenshtein/levenshtein"
)

// Application describes an application with subcommand support.
type Application interface {
	// GetName returns the 'name' of the application.
	GetName() string

	// GetTitle returns a one-line title explaining the purpose of the tool.
	GetTitle() string

	// GetCommands returns the list of the subcommands that are supported.
	GetCommands() []*Command

	// GetOut is used for testing to allow parallel test case execution, should
	// be normally os.Stdout.
	GetOut() io.Writer

	// GetOut is used for testing to allow parallel test case execution, should
	// be normally os.Stderr.
	GetErr() io.Writer

	// GetEnvVars returns the map of EnvVarName -> EnvVarDefinition that this
	// Application responds to.
	GetEnvVars() map[string]EnvVarDefinition
}

// EnvVarDefinition describes an environment variable that this application
// responds to.
type EnvVarDefinition struct {
	Advanced  bool
	ShortDesc string
	Default   string
}

// DefaultApplication implements all of Application interface's methods. An
// application should usually have a global instance of DefaultApplication and
// route main() to command_support.Run(app).
type DefaultApplication struct {
	Name     string
	Title    string
	Commands []*Command
	EnvVars  map[string]EnvVarDefinition
}

// GetName implements interface Application.
func (a *DefaultApplication) GetName() string {
	return a.Name
}

// GetTitle implements interface Application.
func (a *DefaultApplication) GetTitle() string {
	return a.Title
}

// GetCommands implements interface Application.
func (a *DefaultApplication) GetCommands() []*Command {
	return a.Commands
}

// GetOut implements interface Application.
func (a *DefaultApplication) GetOut() io.Writer {
	return os.Stdout
}

// GetErr implements interface Application.
func (a *DefaultApplication) GetErr() io.Writer {
	return os.Stderr
}

// GetEnvVars implements interface Application.
func (a *DefaultApplication) GetEnvVars() map[string]EnvVarDefinition {
	return a.EnvVars
}

// Env is the mapping of resolved environment variables passed to
// CommandRun.Run.
type Env map[string]EnvVar

// EnvVar will document the value and existance of a given environment variable,
// as defined by Application.GetEnvVars. Value will be the value from the
// environment, or the Default value if it didn't exist. Exists will be true iff
// the value was present in the environment.
type EnvVar struct {
	Value  string
	Exists bool
}

// CommandRun is an initialized object representing a subcommand that is ready
// to be executed.
type CommandRun interface {
	// Run execute the actual command. When this function is called by
	// command_support.Run(), the flags have already been parsed.
	Run(a Application, args []string, env Env) int

	// GetFlags returns the flags for this specific command.
	GetFlags() *flag.FlagSet
}

// CommandRunBase implements GetFlags of CommandRun. It should be embedded in
// another struct that implements Run().
type CommandRunBase struct {
	Flags flag.FlagSet
}

// GetFlags implements CommandRun.
func (c *CommandRunBase) GetFlags() *flag.FlagSet {
	return &c.Flags
}

// Command describes a subcommand. It has one generator to generate a command
// object which is executable. The purpose of this design is to enable safe
// parallel execution of test cases.
type Command struct {
	UsageLine  string
	ShortDesc  string
	LongDesc   string
	Advanced   bool
	CommandRun func() CommandRun
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

// usage prints out the general application usage.
//
// TODO(maruel): Use termbox-go to enable coloring!
func usage(out io.Writer, a Application, includeAdvanced bool) {
	usageTemplate := `{{.Title}}

Usage:  {{.Name}} [command] [arguments]

Commands:{{range .Commands}}
  {{.Name | printf "%%-%ds"}}  {{.ShortDesc}}{{end}}

{{if .EnvVars}}Environment Variables:{{range .EnvVars}}
  {{.Name | printf "%%-%ds"}}  {{.ShortDesc}}{{if .Default}} (Default: {{.Default | printf "%%q"}}){{end}}{{end}}

{{end}}
Use "{{.Name}} help [command]" for more information about a command.{{if .ShowAdvancedTip}}
Use "{{.Name}} help -advanced" to display all commands.{{end}}

`

	widestCmd := 0
	allCmds := a.GetCommands()
	cmds := make([]*Command, 0, len(allCmds))
	hasAdvanced := false
	for _, c := range allCmds {
		hasAdvanced = hasAdvanced || c.Advanced

		if !c.Advanced || includeAdvanced {
			// We need to include this command
			if namLen := len(c.Name()); namLen > widestCmd {
				widestCmd = namLen
			}
			cmds = append(cmds, c)
		}
	}

	type envVarEntry struct {
		Name      string
		ShortDesc string
		Default   string
	}
	widestEnvVar := 0
	envVars := []envVarEntry(nil)
	if envVarMap := a.GetEnvVars(); len(envVarMap) > 0 {
		envVarKeys := make(sort.StringSlice, 0, len(envVarMap))
		for k, v := range envVarMap {
			if v.Advanced {
				hasAdvanced = true
			}
			if !v.Advanced || includeAdvanced {
				if keyLen := len(k); keyLen > widestEnvVar {
					widestEnvVar = keyLen
				}
				envVarKeys = append(envVarKeys, k)
			}
		}
		envVarKeys.Sort()
		envVars = make([]envVarEntry, 0, len(envVarKeys))
		for _, k := range envVarKeys {
			v := envVarMap[k]
			envVars = append(envVars, envVarEntry{k, v.ShortDesc, v.Default})
		}
	}
	data := map[string]interface{}{
		"Title":           a.GetTitle(),
		"Name":            a.GetName(),
		"Commands":        cmds,
		"EnvVars":         envVars,
		"ShowAdvancedTip": (hasAdvanced && !includeAdvanced),
	}
	tmpl(out, fmt.Sprintf(usageTemplate, widestCmd, widestEnvVar), data)
}

func getCommandUsageHandler(out io.Writer, a Application, c *Command, r CommandRun, helpUsed *bool) func() {
	return func() {
		helpTemplate := "{{.Cmd.LongDesc | trim | wrapWithLines}}usage:  {{.App.GetName}} {{.Cmd.UsageLine}}\n"
		dict := struct {
			App Application
			Cmd *Command
		}{a, c}
		tmpl(out, helpTemplate, dict)
		r.GetFlags().PrintDefaults()
		*helpUsed = true
	}
}

// Initializes the flags for a specific CommandRun.
func initCommand(a Application, c *Command, r CommandRun, out io.Writer, helpUsed *bool) {
	r.GetFlags().Usage = getCommandUsageHandler(out, a, c, r, helpUsed)
	r.GetFlags().SetOutput(out)
	r.GetFlags().Init(c.Name(), flag.ContinueOnError)
}

// FindCommand finds a Command by name and returns it if found.
func FindCommand(a Application, name string) *Command {
	for _, c := range a.GetCommands() {
		if c.Name() == name {
			return c
		}
	}
	return nil
}

// FindNearestCommand heuristically finds a Command the user wanted to type but
// failed to type correctly.
func FindNearestCommand(a Application, name string) *Command {
	commands := map[string]*Command{}
	for _, c := range a.GetCommands() {
		commands[c.Name()] = c
	}
	if c, ok := commands[name]; ok {
		return c
	}

	// Search for unique prefix.
	withPrefix := []*Command{}
	for n, c := range commands {
		if strings.HasPrefix(n, name) {
			withPrefix = append(withPrefix, c)
		}
	}
	if len(withPrefix) == 1 {
		return withPrefix[0]
	}

	// Search for case insensitivity.
	withPrefix = []*Command{}
	lowName := strings.ToLower(name)
	for n, c := range commands {
		if strings.HasPrefix(strings.ToLower(n), lowName) {
			withPrefix = append(withPrefix, c)
		}
	}
	if len(withPrefix) == 1 {
		return withPrefix[0]
	}

	// Calculate the levenshtein distance and take the closest one.
	closestD := 1000
	var closestC *Command
	secondD := 1000
	for n, c := range commands {
		dist := levenshtein.DistanceForStrings([]rune(n), []rune(name), levenshtein.DefaultOptions)
		if dist < closestD {
			secondD = closestD
			closestD = dist
			closestC = c
		} else if dist < secondD {
			secondD = dist
		}
	}
	if closestD > 3 {
		// Not similar enough. Don't be a fool and run a random command.
		return nil
	}
	if (secondD - closestD) < 3 {
		// Too ambiguous.
		return nil
	}
	return closestC
}

// Run runs the application, scheduling the subcommand. This is the main entry
// point of the library.
func Run(a Application, args []string) int {
	var helpUsed bool

	// Process general flags first, mainly for -help.
	flag.Usage = func() {
		usage(a.GetErr(), a, false)
		helpUsed = true
	}

	// Do not parse during unit tests because flag.commandLine.errorHandling == ExitOnError. :(
	// It is safer to use a base class embedding CommandRunBase that is then
	// embedded by each CommandRun implementation to define flags available for
	// all commands.
	if args == nil {
		flag.Parse()
		args = flag.Args()
	}

	if len(args) < 1 {
		// Need a command.
		usage(a.GetErr(), a, false)
		return 2
	}

	if c := FindNearestCommand(a, args[0]); c != nil {
		// Initialize the flags.
		r := c.CommandRun()
		initCommand(a, c, r, a.GetErr(), &helpUsed)
		if err := r.GetFlags().Parse(args[1:]); err != nil {
			return 2
		}
		if helpUsed {
			return 0
		}
		envVars := a.GetEnvVars()
		envMap := make(map[string]EnvVar, len(envVars))
		for k, v := range envVars {
			val, ok := os.LookupEnv(k)
			if !ok {
				val = v.Default
			}
			envMap[k] = EnvVar{val, ok}
		}
		return r.Run(a, r.GetFlags().Args(), envMap)
	}

	fmt.Fprintf(a.GetErr(), "%s: unknown command %#q\n\nRun '%s help' for usage.\n", a.GetName(), args[0], a.GetName())
	return 2
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace, "wrapWithLines": wrapWithLines})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(fmt.Sprintf("Failed to execute templare: %s", err))
	}
}

func wrapWithLines(s string) string {
	if s == "" {
		return s
	}
	return s + "\n\n"
}

// CmdHelp defines the help command. It should be included in your application's
// Commands list.
//
// It is not added automatically but it will be run automatically if added.
var CmdHelp = &Command{
	UsageLine: "help [<command>|-advanced]",
	ShortDesc: "prints help about a command",
	LongDesc:  "Prints an overview of every command or information about a specific command.\nPass -advanced to see help for advanced commands.",
	CommandRun: func() CommandRun {
		ret := &helpRun{}
		ret.Flags.BoolVar(&ret.advanced, "advanced", false, "show advanced commands")
		return ret
	},
}

type helpRun struct {
	CommandRunBase
	advanced bool
}

func (c *helpRun) Run(a Application, args []string, env Env) int {
	if len(args) == 0 {
		usage(a.GetOut(), a, c.advanced)
		return 0
	}
	if len(args) != 1 {
		fmt.Fprintf(a.GetErr(), "%s: Too many arguments given\n\nRun '%s help' for usage.\n", a.GetName(), a.GetName())
		return 2
	}
	// Redirects all output to Out.
	var helpUsed bool
	if cmd := FindNearestCommand(a, args[0]); cmd != nil {
		// Initialize the flags.
		r := cmd.CommandRun()
		initCommand(a, cmd, r, a.GetErr(), &helpUsed)
		r.GetFlags().Usage()
		return 0
	}

	fmt.Fprintf(a.GetErr(), "%s: unknown command %#q\n\nRun '%s help' for usage.\n", a.GetName(), args[0], a.GetName())
	return 2
}

package cli

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

// App is the main structure of a cli application. It is recomended that
// and app be created with the cli.NewApp() function.
type App struct {
	// The name of the program. Defaults to os.Args[0]
	Name string
	// The program executable. Defaults to os.Args[0]
	Exec string
	// Description of the program
	Description string
	// Usage pattern for executing the program
	Usage string
	// Version of the program
	Version string
	// List of commands to execute
	Commands []Command
	// List of flags to parse
	Flags []Flag
	// An action to execute before any commands are run, but after the context is ready
	// If a non-nil error is returned, no commands are run
	Before func(context *Context) error
	// The action to execute when no command is specified
	Action func(context *Context)
	// Execute this function if the proper command cannot be found
	CommandNotFound func(context *Context, command string)
	// Compilation date
	Compiled time.Time
	// Author
	Author string
	// Author e-mail
	Email string
}

// NewApp creates a new cli Application with some reasonable defaults.
func NewApp() *App {
	return &App{
		Name:        os.Args[0],
		Exec:        os.Args[0],
		Description: "A new application",
		Usage:       os.Args[0] + " [options] <command>",
		Version:     "0.0.0",
		Action:      helpCommand.Action,
		Compiled:    compileTime(),
	}
}

// Run is the entry point to the cli app. It parses the arguments slice and routes to the
// proper flag/args combination.
func (a *App) Run(arguments []string) error {
	// append help to commands
	if a.Command(helpCommand.Name) == nil {
		a.Commands = append(a.Commands, helpCommand)
	}

	// append version flag
	a.appendFlag(VersionFlag)

	// parse flags
	set := flagSet(a.Name, a.Flags)
	set.SetOutput(ioutil.Discard)
	err := set.Parse(arguments[1:])
	nerr := normalizeFlags(a.Flags, set)
	if nerr != nil {
		fmt.Println(nerr)
		context := NewContext(a, set, set)
		ShowAppHelp(context)
		fmt.Println("")
		return nerr
	}
	context := NewContext(a, set, set)

	if err != nil {
		fmt.Printf("Incorrect Usage - type '%s help' for info\n\n", a.Exec)
		return err
	}

	if checkVersion(context) {
		return nil
	}

	if a.Before != nil {
		err := a.Before(context)
		if err != nil {
			return err
		}
	}

	args := context.Args()
	if args.Present() {
		name := args.First()
		c := a.Command(name)
		if c != nil {
			return c.Run(context)
		}
	}

	// Run default Action
	a.Action(context)
	return nil
}

// RunAndExitOnError is another entry point to the cli app. It takes care of passing
// arguments and error handling.
func (a *App) RunAndExitOnError() {
	if err := a.Run(os.Args); err != nil {
		os.Stderr.WriteString(fmt.Sprintln(err))
		os.Exit(1)
	}
}

// Command returns the named command on App. Returns nil if the command does not exist.
func (a *App) Command(name string) *Command {
	for _, c := range a.Commands {
		if c.HasName(name) {
			return &c
		}
	}

	return nil
}

func (a *App) hasFlag(flag Flag) bool {
	for _, f := range a.Flags {
		if flag == f {
			return true
		}
	}

	return false
}

func (a *App) appendFlag(flag Flag) {
	if !a.hasFlag(flag) {
		a.Flags = append(a.Flags, flag)
	}
}

// Tries to find out when this binary was compiled.
// Returns the current time if it fails to find it.
func compileTime() time.Time {
	info, err := os.Stat(os.Args[0])
	if err != nil {
		return time.Now()
	}
	return info.ModTime()
}

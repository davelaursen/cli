package cli

import (
	"fmt"
	"io/ioutil"
)

// Command is a subcommand for a cli.App.
type Command struct {
	// The name of the command
	Name string
	// A short description of the command
	ShortDescription string
	// Usage pattern for executing the command
	Usage string
	// A longer explanation of how the command works
	Description string
	// An action to execute before any subcommands are run, but after the context is ready
	// If a non-nil error is returned, no subcommands are run
	Before func(context *Context) error
	// The function to call when this command is invoked without a subcommand
	Action func(context *Context)
	// List of child commands
	Subcommands []Subcommand
	// List of flags to parse
	Flags []Flag
}

// Run invokes the command given the context. It parses ctx.Args() to generate command-specific flags.
func (c Command) Run(ctx *Context) error {
	set := flagSet(c.Name, c.Flags)
	set.SetOutput(ioutil.Discard)
	err := set.Parse(ctx.Args()[1:])
	if err != nil {
		fmt.Printf("Incorrect Usage - type '%s help' for info\n\n", ctx.App.Exec)
		return err
	}

	nerr := normalizeFlags(c.Flags, set)
	if nerr != nil {
		fmt.Println(nerr)
		fmt.Println("")
		ShowCommandHelp(ctx, c.Name)
		fmt.Println("")
		return nerr
	}

	context := NewContext(ctx.App, set, ctx.globalSet)
	if len(c.Subcommands) > 0 && len(set.Args()) > 0 {
		name := set.Args()[0]
		s := c.Subcommand(name)
		if s != nil {
			return s.Run(context)
		}
	}

	context.Command = c
	c.Action(context)
	return nil
}

// HasName returns true if Command.Name matches given name
func (c Command) HasName(name string) bool {
	return c.Name == name
}

// Command returns the named command on App. Returns nil if the command does not exist.
func (c Command) Subcommand(name string) *Subcommand {
	for _, s := range c.Subcommands {
		if s.HasName(name) {
			return &s
		}
	}

	return nil
}

// Command is a subcommand for a cli.App.
type Subcommand struct {
	// The name of the command
	Name string
	// A short description of the command
	ShortDescription string
	// Usage pattern for executing the command
	Usage string
	// A longer explanation of how the command works
	Description string
	// The function to call when this command is invoked without a subcommand
	Action func(context *Context)
	// List of flags to parse
	Flags []Flag
}

// Run invokes the command given the context. It parses ctx.Args() to generate command-specific flags.
func (s Subcommand) Run(ctx *Context) error {
	set := flagSet(s.Name, s.Flags)
	set.SetOutput(ioutil.Discard)
	err := set.Parse(ctx.Args()[1:])
	if err != nil {
		fmt.Printf("Incorrect Usage - type '%s help' for info\n\n", ctx.App.Exec)
		return err
	}

	nerr := normalizeFlags(s.Flags, set)
	if nerr != nil {
		fmt.Println(nerr)
		fmt.Println("")
		ShowCommandHelp(ctx, s.Name)
		fmt.Println("")
		return nerr
	}

	context := NewContext(ctx.App, set, ctx.globalSet)
	context.Command = subcmdToCmd(s)
	s.Action(context)
	return nil
}

// HasName returns true if Subcommand.Name matches given name
func (s Subcommand) HasName(name string) bool {
	return s.Name == name
}

func subcmdToCmd(s Subcommand) Command {
	return Command{
		Name:             s.Name,
		ShortDescription: s.ShortDescription,
		Usage:            s.Usage,
		Description:      s.Description,
		Action:           s.Action,
		Flags:            s.Flags,
	}
}

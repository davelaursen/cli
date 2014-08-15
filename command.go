package cli

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// Command is a subcommand for a cli.App.
type Command struct {
	// The name of the command
	Name string
	// A short description of the usage of this command
	Summary string
	// Usage pattern
	Usage string
	// A longer explanation of how the command works
	Description string
	// The function to call when this command is invoked
	Action func(context *Context)
	// List of flags to parse
	Flags []Flag
}

// Run invokes the command given the context, parses ctx.Args() to generate command-specific flags.
func (c Command) Run(ctx *Context) error {
	set := flagSet(c.Name, c.Flags)
	set.SetOutput(ioutil.Discard)

	firstFlagIndex := -1
	for index, arg := range ctx.Args() {
		if strings.HasPrefix(arg, "-") {
			firstFlagIndex = index
			break
		}
	}

	var err error
	if firstFlagIndex > -1 {
		args := ctx.Args()
		regularArgs := args[1:firstFlagIndex]
		flagArgs := args[firstFlagIndex:]
		err = set.Parse(append(flagArgs, regularArgs...))
	} else {
		err = set.Parse(ctx.Args().Tail())
	}

	if err != nil {
		fmt.Printf("Incorrect Usage - type '%s help' for info\n\n", ctx.App.Exec)
		return err
	}

	nerr := normalizeFlags(c.Flags, set)
	if nerr != nil {
		fmt.Printf("Incorrect Usage - type '%s help' for info\n\n", ctx.App.Exec)
		fmt.Print(nerr, "\n\n")
		return nerr
	}
	context := NewContext(ctx.App, set, ctx.globalSet)

	if checkCommandHelp(context, c.Name) {
		return nil
	}
	context.Command = c
	c.Action(context)
	return nil
}

// HasName returns true if Command.Name matches given name
func (c Command) HasName(name string) bool {
	return c.Name == name
}

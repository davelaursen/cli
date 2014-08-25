package cli

import (
	"fmt"
	"os"
	"testing"
)

func ExampleApp() {
	// set args for examples sake
	os.Args = []string{"greet", "--name", "Jeremy"}

	app := NewApp()
	app.Name = "greet"
	app.Flags = []Flag{
		StringFlag{Name: "name", Value: "bob", Description: "a name to say"},
	}
	app.Action = func(c *Context) {
		fmt.Printf("Hello %v\n", c.String("name"))
	}
	app.Run(os.Args)
	// Output:
	// Hello Jeremy
}

func ExampleAppHelp() {
	// set args for examples sake
	os.Args = []string{"greet", "help"}

	app := NewApp()
	app.Name = "greet"
	app.Flags = []Flag{
		StringFlag{Name: "name", Value: "bob", Description: "a name to say"},
	}
	app.Commands = []Command{
		{
			Name:             "describeit",
			ShortDescription: "use it to see a description",
			Usage:            "test",
			Description:      "This is how we describe describeit the function",
			Action: func(c *Context) {
				fmt.Printf("i like to describe things")
			},
		},
	}
	app.Run(os.Args)
	// Output:
	// greet, v0.0.0
	// A new application
	//
	// USAGE:
	//    greet [options] <command>
	//
	// COMMANDS:
	//    describeit	use it to see a description
	//    help		Shows a list of commands or help for one command
	//
	//    Use 'greet help <command> [<subcommand>]' for more
	//    information about a command or subcommand.
	//
	// OPTIONS:
	//    -name 'bob'	a name to say
	//    -version	print the version
}

func TestApp_Run(t *testing.T) {
	s := ""

	app := NewApp()
	app.Action = func(c *Context) {
		s = s + c.Args().First()
	}

	err := app.Run([]string{"command", "foo"})
	expect(t, err, nil)
	err = app.Run([]string{"command", "bar"})
	expect(t, err, nil)
	expect(t, s, "foobar")
}

var commandAppTests = []struct {
	name     string
	expected bool
}{
	{"foobar", true},
	{"batbaz", true},
	{"b", false},
	{"f", false},
	{"bat", false},
	{"nothing", false},
}

func TestApp_Command(t *testing.T) {
	app := NewApp()
	fooCommand := Command{Name: "foobar"}
	batCommand := Command{Name: "batbaz"}
	app.Commands = []Command{
		fooCommand,
		batCommand,
	}

	for _, test := range commandAppTests {
		expect(t, app.Command(test.name) != nil, test.expected)
	}
}

func TestApp_Float64Flag(t *testing.T) {
	var meters float64

	app := NewApp()
	app.Flags = []Flag{
		Float64Flag{Name: "height", Value: 1.5, Description: "Set the height, in meters"},
	}
	app.Action = func(c *Context) {
		meters = c.Float64("height")
	}

	app.Run([]string{"", "--height", "1.93"})
	expect(t, meters, 1.93)
}

func TestAppHelpPrinter(t *testing.T) {
	oldPrinter := HelpPrinter
	defer func() {
		HelpPrinter = oldPrinter
	}()

	var wasCalled = false
	HelpPrinter = func(template string, data interface{}) {
		wasCalled = true
	}

	app := NewApp()
	app.Run([]string{"-h"})

	if wasCalled == false {
		t.Errorf("Help printer expected to be called, but was not")
	}
}

func TestAppVersionPrinter(t *testing.T) {
	oldPrinter := VersionPrinter
	defer func() {
		VersionPrinter = oldPrinter
	}()

	var wasCalled = false
	VersionPrinter = func(c *Context) {
		wasCalled = true
	}

	app := NewApp()
	ctx := NewContext(app, nil, nil)
	ShowVersion(ctx)

	if wasCalled == false {
		t.Errorf("Version printer expected to be called, but was not")
	}
}

func TestAppCommandNotFound(t *testing.T) {
	beforeRun, subcommandRun := false, false
	app := NewApp()

	app.CommandNotFound = func(c *Context, command string) {
		beforeRun = true
	}

	app.Commands = []Command{
		Command{
			Name: "bar",
			Action: func(c *Context) {
				subcommandRun = true
			},
		},
	}

	app.Run([]string{"command", "foo"})

	expect(t, beforeRun, true)
	expect(t, subcommandRun, false)
}

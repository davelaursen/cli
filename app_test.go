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
		StringFlag{Name: "name", Value: "bob", Usage: "a name to say"},
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
		StringFlag{Name: "name", Value: "bob", Usage: "a name to say"},
	}
	app.Commands = []Command{
		{
			Name:        "describeit",
			Summary:     "use it to see a description",
			Usage:       "test",
			Description: "This is how we describe describeit the function",
			Action: func(c *Context) {
				fmt.Printf("i like to describe things")
			},
		},
	}
	app.Run(os.Args)
	// greet, v0.0.0 - A new application
	//
	// USAGE:
	//    greet [global options] command [command options] [arguments...]
	//
	// COMMANDS:
	//    describeit	use it to see a description
	//    help		Shows a list of commands or help for one command
	//
	// OPTIONS:
	//    --name 'bob'	a name to say
	//    --version	print the version
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

func TestApp_CommandWithArgBeforeFlags(t *testing.T) {
	var parsedOption, firstArg string

	app := NewApp()
	command := Command{
		Name: "cmd",
		Flags: []Flag{
			StringFlag{Name: "option", Value: "", Usage: "some option"},
		},
		Action: func(c *Context) {
			parsedOption = c.String("option")
			firstArg = c.Args().First()
		},
	}
	app.Commands = []Command{command}

	app.Run([]string{"", "cmd", "my-arg", "--option", "my-option"})

	expect(t, parsedOption, "my-option")
	expect(t, firstArg, "my-arg")
}

func TestApp_Float64Flag(t *testing.T) {
	var meters float64

	app := NewApp()
	app.Flags = []Flag{
		Float64Flag{Name: "height", Value: 1.5, Usage: "Set the height, in meters"},
	}
	app.Action = func(c *Context) {
		meters = c.Float64("height")
	}

	app.Run([]string{"", "--height", "1.93"})
	expect(t, meters, 1.93)
}

func TestApp_ParseSliceFlags(t *testing.T) {
	var parsedOption, firstArg string
	var parsedIntSlice []int
	var parsedStringSlice []string

	app := NewApp()
	command := Command{
		Name: "cmd",
		Flags: []Flag{
			IntSliceFlag{Name: "p", Value: &IntSlice{}, Usage: "set one or more ip addr"},
			StringSliceFlag{Name: "ip", Value: &StringSlice{}, Usage: "set one or more ports to open"},
		},
		Action: func(c *Context) {
			parsedIntSlice = c.IntSlice("p")
			parsedStringSlice = c.StringSlice("ip")
			parsedOption = c.String("option")
			firstArg = c.Args().First()
		},
	}
	app.Commands = []Command{command}

	app.Run([]string{"", "cmd", "my-arg", "-p", "22", "-p", "80", "-ip", "8.8.8.8", "-ip", "8.8.4.4"})

	IntsEquals := func(a, b []int) bool {
		if len(a) != len(b) {
			return false
		}
		for i, v := range a {
			if v != b[i] {
				return false
			}
		}
		return true
	}

	StrsEquals := func(a, b []string) bool {
		if len(a) != len(b) {
			return false
		}
		for i, v := range a {
			if v != b[i] {
				return false
			}
		}
		return true
	}
	var expectedIntSlice = []int{22, 80}
	var expectedStringSlice = []string{"8.8.8.8", "8.8.4.4"}

	if !IntsEquals(parsedIntSlice, expectedIntSlice) {
		t.Errorf("%v does not match %v", parsedIntSlice, expectedIntSlice)
	}

	if !StrsEquals(parsedStringSlice, expectedStringSlice) {
		t.Errorf("%v does not match %v", parsedStringSlice, expectedStringSlice)
	}
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

package cli

import (
	"flag"
	"testing"
)

func TestCommand(t *testing.T) {
	app := NewApp()
	set := flag.NewFlagSet("test", 0)
	test := []string{"blah", "-break"}
	set.Parse(test)

	c := NewContext(app, set, set)

	command := Command{
		Name:             "test-cmd",
		ShortDescription: "this is for testing",
		Usage:            "test",
		Description:      "testing",
		Action:           func(_ *Context) {},
	}
	err := command.Run(c)

	expect(t, err.Error(), "flag provided but not defined: -break")
}

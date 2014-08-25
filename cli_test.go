package cli

import (
	"os"
)

func Example() {
	app := NewApp()
	app.Name = "todo"
	app.Usage = "task list on the command line"
	app.Commands = []Command{
		{
			Name:             "add",
			ShortDescription: "add a task to the list",
			Usage:            "add",
			Action: func(c *Context) {
				println("added task: ", c.Args().First())
			},
		},
		{
			Name:             "complete",
			ShortDescription: "complete a task on the list",
			Usage:            "complete",
			Action: func(c *Context) {
				println("completed task: ", c.Args().First())
			},
		},
	}

	app.Run(os.Args)
}

func ExampleSubcommand() {
	app := NewApp()
	app.Name = "say"
	app.Commands = []Command{
		{
			Name:             "hello",
			ShortDescription: "use it to see a description",
			Usage:            "hello",
			Description:      "This is how we describe hello the function",
		}, {
			Name:  "bye",
			Usage: "says goodbye",
			Action: func(c *Context) {
				println("bye")
			},
		},
	}

	app.Run(os.Args)
}

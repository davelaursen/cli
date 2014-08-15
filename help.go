package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"text/template"
)

// The text template for the Default help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var AppHelpTemplate = `
   {{.Name}}, v{{.Version}} - {{.Summary}}

USAGE:
   {{.Usage}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{ "\t" }}{{.Summary}}
   {{end}}{{ if .Flags }}
OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

// The text template for the command help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var CommandHelpTemplate = `
   {{.Name}} - {{.Summary}}

USAGE:
   {{.Usage}}

DESCRIPTION:
   {{.Description}}{{ if .Flags }}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{ end }}
`

var helpCommand = Command{
	Name:    "help",
	Summary: "Shows a list of commands or help for one command",
	Usage:   "help [command]",
	Action: func(c *Context) {
		args := c.Args()
		if args.Present() {
			ShowCommandHelp(c, args.First())
		} else {
			ShowAppHelp(c)
		}
	},
}

// Prints help for the App
var HelpPrinter = printHelp

// Prints version for the App
var VersionPrinter = printVersion

func ShowAppHelp(c *Context) {
	HelpPrinter(AppHelpTemplate, c.App)
}

// Prints help for the given command
func ShowCommandHelp(c *Context, command string) {
	for _, c := range c.App.Commands {
		if c.HasName(command) {
			HelpPrinter(CommandHelpTemplate, c)
			return
		}
	}

	if c.App.CommandNotFound != nil {
		c.App.CommandNotFound(c, command)
	} else {
		fmt.Printf("No help topic for '%v'\n", command)
	}
}

// Prints the version number of the App
func ShowVersion(c *Context) {
	VersionPrinter(c)
}

func printVersion(c *Context) {
	fmt.Printf("%v version %v\n", c.App.Name, c.App.Version)
}

func printHelp(templ string, data interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	t := template.Must(template.New("help").Parse(templ))
	err := t.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func checkVersion(c *Context) bool {
	if c.GlobalBool("version") {
		ShowVersion(c)
		return true
	}

	return false
}

func checkHelp(c *Context) bool {
	if c.GlobalBool("h") || c.GlobalBool("help") {
		ShowAppHelp(c)
		return true
	}

	return false
}

func checkCommandHelp(c *Context, name string) bool {
	if c.Bool("h") || c.Bool("help") {
		ShowCommandHelp(c, name)
		return true
	}

	return false
}

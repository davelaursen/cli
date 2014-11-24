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
{{.Name}}, v{{.Version}}
{{.Description}}

USAGE:
   {{.Usage}}

COMMANDS:
{{range .Commands}}{{ "   " }}{{.Name}}{{ "\t" }}{{.ShortDescription}}{{ "\n" }}{{end}}
   Use '{{.Exec}} help <command> [<subcommand>]' for more
   information about a command or subcommand.
{{ if .Flags }}
OPTIONS:
{{range .Flags}}{{ "   " }}{{.}}{{ "\n" }}{{end}}{{ end }}
`

// The text template for the command help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var CommandHelpTemplate = `
{{.Name}} - {{.ShortDescription}}

USAGE:
   {{.Usage}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}{{if .Subcommands}}

SUBCOMMANDS:
   {{range .Subcommands}}{{.Name}}{{ "\t" }}{{.ShortDescription}}
   {{end}}{{end}}
`

// The text template for the subcommand help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var SubcommandHelpTemplate = `
{{.Name}} - {{.ShortDescription}}

USAGE:
   {{.Usage}}{{if .Description}}

DESCRIPTION:
   {{.Description}}{{end}}
`

var helpCommand = Command{
	Name:             "help",
	ShortDescription: "Shows a list of commands or help for one command",
	Action: func(c *Context) {
		args := c.Args()
		if args.Present() {
			if len(args) == 2 {
				ShowSubcommandHelp(c, args[0], args[1])
			} else {
				ShowCommandHelp(c, args.First())
			}
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

// Prints the list of subcommands as the default app completion method
func DefaultAppComplete(c *Context) {
	for _, command := range c.App.Commands {
		fmt.Println(command.Name)
	}
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

// Prints help for the given subcommand
func ShowSubcommandHelp(ctx *Context, command, subcommand string) {
	for _, c := range ctx.App.Commands {
		if c.HasName(command) {
			for _, s := range c.Subcommands {
				if s.HasName(subcommand) {
					HelpPrinter(SubcommandHelpTemplate, s)
					return
				}
			}
		}
	}

	if ctx.App.CommandNotFound != nil {
		ctx.App.CommandNotFound(ctx, command)
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

# cli.go
cli.go is a simple, fast, and fun package for building command line apps in Go. The goal is to enable developers to write fast and distributable command line applications in an expressive way.

Note: This package provides a subset of the functionality provided by the original [cli.go](https://github.com/codegangsta/cli) package created by Codegangsta. It simplies the API, reformats the help output, and adds some minor functionality.

## Overview
Command line apps are usually so tiny that there is absolutely no reason why your code should *not* be self-documenting. Things like generating help text and parsing command flags/options should not hinder productivity when writing a command line app.

This is where cli.go comes into play. cli.go makes command line programming fun, organized, and expressive!

## Installation
Make sure you have a working Go environment (go 1.1 is *required*). [See the install instructions](http://golang.org/doc/install.html).

To install cli.go, simply run:
```
$ go get github.com/davelaursen/cli
```

Make sure your PATH includes the `$GOPATH/bin` directory so your commands can be easily used:
```
export PATH=$PATH:$GOPATH/bin
```

## Getting Started
One of the philosophies behind cli.go is that an API should be playful and full of discovery. So a cli.go app can be as little as one line of code in `main()`. 

``` go
package main

import (
  "os"
  "github.com/davelaursen/cli"
)

func main() {
  cli.NewApp().Run(os.Args)
}
```

This app will run and show help text, but is not very useful. Let's give an action to execute and some help documentation:

``` go
package main

import (
  "os"
  "github.com/davelaursen/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "boom"
  app.Summary = "make an explosive entrance"
  app.Usage = "boom [help]"
  app.Action = func(c *cli.Context) {
    println("boom! I say!")
  }
  
  app.Run(os.Args)
}
```

Running this already gives you a ton of functionality, plus support for things like flags, which are covered below.

## Example

Being a programmer can be a lonely job. Thankfully by the power of automation that is not the case! Let's create a greeter app to fend off our demons of loneliness!

``` go
/* greet.go */
package main

import (
  "os"
  "github.com/davelaursen/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "greet"
  app.Summary = "fight the loneliness!" 
  app.Usage = "greet [help]"
  app.Action = func(c *cli.Context) {
    println("Hello friend!")
  }
  
  app.Run(os.Args)
}
```

Install our command to the `$GOPATH/bin` directory:

```
$ go install
```

Finally run our new command:

```
$ greet
Hello friend!
```

cli.go also generates some bitchass help text:
```
$ greet help
greet, v0.0.0 - fight the loneliness!

USAGE:
   greet [help]

COMMANDS:
   help       Shows a list of commands or help for one command

OPTIONS:
   --version  print the version
```

### Arguments
You can lookup arguments by calling the `Args` function on cli.Context.

``` go
...
app.Action = func(c *cli.Context) {
  println("Hello", c.Args()[0])
}
...
```

### Flags
Setting and querying flags is simple.
``` go
...
app.Flags = []cli.Flag {
  cli.StringFlag{
    Name: "lang",
    Value: "english",
    Usage: "language for the greeting",
  },
}
app.Action = func(c *cli.Context) {
  name := "someone"
  if len(c.Args()) > 0 {
    name = c.Args()[0]
  }
  if c.String("lang") == "spanish" {
    println("Hola", name)
  } else {
    println("Hello", name)
  }
}
...
```

#### Alternate Names

You can set alternate (or short) names for flags by providing a comma-delimited list for the Name. e.g.

``` go
app.Flags = []cli.Flag {
  cli.StringFlag{
    Name: "lang, l",
    Value: "english",
    Usage: "language for the greeting",
  },
}
```

#### Values from the Environment

You can also have the default value set from the environment via EnvVar.  e.g.

``` go
app.Flags = []cli.Flag {
  cli.StringFlag{
    Name: "lang, l",
    Value: "english",
    Usage: "language for the greeting",
    EnvVar: "APP_LANG",
  },
}
```

That flag can then be set with `--lang spanish` or `-l spanish`. Note that giving two different forms of the same flag in the same command invocation is an error.


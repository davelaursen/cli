package cli

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// This flag enables bash-completion for all commands and subcommands
var BashCompletionFlag = BoolFlag{
	Name: "generate-bash-completion",
}

// This flag prints the version for the application
var VersionFlag = BoolFlag{
	Name:        "version",
	Description: "print the version",
}

// Flag is a common interface related to parsing flags in cli.
// For more advanced flag parsing techniques, it is recomended that
// this interface be implemented.
type Flag interface {
	fmt.Stringer
	// Apply Flag settings to the given flag set
	Apply(*flag.FlagSet)
	getName() string
}

func flagSet(name string, flags []Flag) *flag.FlagSet {
	set := flag.NewFlagSet(name, flag.ContinueOnError)

	for _, f := range flags {
		f.Apply(set)
	}
	return set
}

func eachName(longName string, fn func(string)) {
	parts := strings.Split(longName, ",")
	for _, name := range parts {
		name = strings.Trim(name, " ")
		fn(name)
	}
}

// Generic is a generic parseable type identified by a specific flag
type Generic interface {
	Set(value string) error
	String() string
}

// GenericFlag is the flag type for types implementing Generic
type GenericFlag struct {
	Name        string
	Value       Generic
	Description string
	EnvVar      string
}

func (f GenericFlag) String() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("-%s %v\t`%v` %s", f.Name, f.Value, "-"+f.Name+" option -"+f.Name+" option", f.Description))
}

func (f GenericFlag) Apply(set *flag.FlagSet) {
	val := f.Value
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			val.Set(envVal)
		}
	}

	eachName(f.Name, func(name string) {
		set.Var(f.Value, name, f.Description)
	})
}

func (f GenericFlag) getName() string {
	return f.Name
}

type BoolFlag struct {
	Name        string
	Description string
	EnvVar      string
}

func (f BoolFlag) String() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s\t%v", prefixedNames(f.Name), f.Description))
}

func (f BoolFlag) Apply(set *flag.FlagSet) {
	val := false
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValBool, err := strconv.ParseBool(envVal)
			if err == nil {
				val = envValBool
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Bool(name, val, f.Description)
	})
}

func (f BoolFlag) getName() string {
	return f.Name
}

type BoolTFlag struct {
	Name        string
	Description string
	EnvVar      string
}

func (f BoolTFlag) String() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s\t%v", prefixedNames(f.Name), f.Description))
}

func (f BoolTFlag) Apply(set *flag.FlagSet) {
	val := true
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValBool, err := strconv.ParseBool(envVal)
			if err == nil {
				val = envValBool
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Bool(name, val, f.Description)
	})
}

func (f BoolTFlag) getName() string {
	return f.Name
}

type StringFlag struct {
	Name        string
	Value       string
	Description string
	EnvVar      string
}

func (f StringFlag) String() string {
	var fmtString string
	fmtString = "%s %v\t%v"

	if len(f.Value) > 0 {
		fmtString = "%s '%v'\t%v"
	} else {
		fmtString = "%s %v\t%v"
	}

	return withEnvHint(f.EnvVar, fmt.Sprintf(fmtString, prefixedNames(f.Name), f.Value, f.Description))
}

func (f StringFlag) Apply(set *flag.FlagSet) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			f.Value = envVal
		}
	}

	eachName(f.Name, func(name string) {
		set.String(name, f.Value, f.Description)
	})
}

func (f StringFlag) getName() string {
	return f.Name
}

type IntFlag struct {
	Name        string
	Value       int
	Description string
	EnvVar      string
}

func (f IntFlag) String() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Description))
}

func (f IntFlag) Apply(set *flag.FlagSet) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValInt, err := strconv.ParseUint(envVal, 10, 64)
			if err == nil {
				f.Value = int(envValInt)
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Int(name, f.Value, f.Description)
	})
}

func (f IntFlag) getName() string {
	return f.Name
}

type DurationFlag struct {
	Name        string
	Value       time.Duration
	Description string
	EnvVar      string
}

func (f DurationFlag) String() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Description))
}

func (f DurationFlag) Apply(set *flag.FlagSet) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValDuration, err := time.ParseDuration(envVal)
			if err == nil {
				f.Value = envValDuration
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Duration(name, f.Value, f.Description)
	})
}

func (f DurationFlag) getName() string {
	return f.Name
}

type Float64Flag struct {
	Name        string
	Value       float64
	Description string
	EnvVar      string
}

func (f Float64Flag) String() string {
	return withEnvHint(f.EnvVar, fmt.Sprintf("%s '%v'\t%v", prefixedNames(f.Name), f.Value, f.Description))
}

func (f Float64Flag) Apply(set *flag.FlagSet) {
	if f.EnvVar != "" {
		if envVal := os.Getenv(f.EnvVar); envVal != "" {
			envValFloat, err := strconv.ParseFloat(envVal, 10)
			if err == nil {
				f.Value = float64(envValFloat)
			}
		}
	}

	eachName(f.Name, func(name string) {
		set.Float64(name, f.Value, f.Description)
	})
}

func (f Float64Flag) getName() string {
	return f.Name
}

func prefixedNames(fullName string) (prefixed string) {
	parts := strings.Split(fullName, ",")
	for i, name := range parts {
		name = strings.Trim(name, " ")
		prefixed += "-" + name
		if i < len(parts)-1 {
			prefixed += ", "
		}
	}
	return
}

func withEnvHint(envVar, str string) string {
	envText := ""
	if envVar != "" {
		envText = fmt.Sprintf(" [$%s]", envVar)
	}
	return str + envText
}

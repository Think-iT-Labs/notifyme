package argparser

import (
	"errors"
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

type Arguments struct {
	Verbose    []bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
	ConfigFile string   `short:"c" long:"config" description:"Config file"`
	UserCmd    []string `positional-arg-name:"COMMAND" positional-args:"true" required:"true"`
}

func Parse(args []string) (Arguments, error) {
	var arguments Arguments
	parser := flags.NewParser(&arguments, flags.Default)

	cmdIndex := getUserCmdIndex(args)
	_, err := parser.ParseArgs(args[:cmdIndex])
	if err != nil {
		return Arguments{}, err
	}

	if len(args[cmdIndex:]) == 0 {
		return Arguments{}, errors.New("No command provided")
	}

	arguments.UserCmd = args[cmdIndex:]

	return arguments, nil
}

func MustParse(args []string) Arguments {
	a, err := Parse(args)
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(2)
		}
	}
	return a
}

func getUserCmdIndex(args []string) int {
	var cmdIndex int
	for _, arg := range args {
		if arg[0] != '-' {
			break
		}
		cmdIndex++
	}
	return cmdIndex
}

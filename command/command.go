package command

import (
	"errors"
	"regexp"
	"strings"
)

type Command struct {
	Command   string
	Arguments []string
}

func ParseCommand(s string) (Command, error) {
	if !strings.HasPrefix(s, "--") && len(s) < 3 {
		return Command{}, errors.New("parseCommand: Not a command")
	}

	// Remove double white spaces
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")

	fields := strings.Fields(s)

	cmd := Command{fields[0][2:], fields[1:]}

	return cmd, nil
}

func IsCommand(c *Command, name string) bool {
	return c.Command == name
}

func VerifyArguments(c *Command, count int) error {
	// TODO: find a way to provide expected limitations for each argument
	if len(c.Arguments) != count {
		return errors.New(c.Command + ": Incorrect command arguments")
	}
	return nil
}

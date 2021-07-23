// Package command Serves to regex and fix command inputs from users. Returns processed string arrays with the command and arguments parts.
package command

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strconv"
	"strings"
)

const prefix string = "."

type Command struct {
	Command   string
	Arguments []string
}

type IntegerArg struct {
	LowerLimit int
	UpperLimit int
}

type RegexArg struct {
	Expression   string
	CaptureGroup int
}

func ParseCommand(s string) (Command, error) {
	if !strings.HasPrefix(s, prefix) && len(s) < len(prefix)+1 {
		return Command{}, errors.New("")
	}

	// Remove double white spaces
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")

	fields := strings.Fields(s)

	cmd := Command{fields[0][len(prefix):], fields[1:]}

	return cmd, nil
}

func ParseMentionToString(s string) string {
	s = strings.Replace(s, "<", "", 1)
	s = strings.Replace(s, ">", "", 1)
	s = strings.Replace(s, "!", "", 1)
	s = strings.Replace(s, "@", "", 1)

	return s
}

func IsCommand(c *Command, name string) bool {
	return c.Command == name
}

func VerifyArguments(c *Command, args ...interface{}) error {
	if len(c.Arguments) != len(args) {
		return errors.New(c.Command + ": Incorrect command arguments")
	}

	for i, arg := range args {
		switch t := arg.(type) {
		case int:
			_, err := strconv.ParseInt(c.Arguments[i], 10, 64)
			if err != nil {
				return printArgError(c.Command, c.Arguments[i], "is not a number")
			}

		case string:
			if t != c.Arguments[i] {
				return printArgError(c.Command, c.Arguments[i], "isn't the expected argument "+t)
			}

		case IntegerArg:
			n, err := strconv.Atoi(c.Arguments[i])
			if err != nil || n < t.LowerLimit || n > t.UpperLimit {
				return printArgError(c.Command, c.Arguments[i], "is not a number between"+strconv.Itoa(t.LowerLimit)+
					" and "+strconv.Itoa(t.UpperLimit))
			}

		case RegexArg:
			re, err := regexp.Compile(t.Expression)
			if err != nil {
				return printError(c.Command, "Internal error. Regex for argument["+strconv.Itoa(i)+"] can't be compiled")
			}

			matches := re.FindStringSubmatch(c.Arguments[i])
			if len(matches) == 0 {
				return printArgError(c.Command, c.Arguments[i], "is not a valid argument")
			}

			// Export desired capture-group
			c.Arguments[i] = matches[t.CaptureGroup]

		default:
			return printError(c.Command, "Internal error")

		}
	}

	return nil
}

func printError(command string, cause string) error {
	return fmt.Errorf("%s: %s", command, cause)
}

func printArgError(command string, argument string, cause string) error {
	return printError(command, fmt.Sprintf("Argument \"%s\" %s", argument, cause))
}

// JoinArguments takes all arguments separated by space and joins them together into a single string
//TODO: make this function take arguments as to which fields of the cmd.Arguments[x] to unify
func JoinArguments(cmd Command) string {
	var commandString string
	for _, value := range cmd.Arguments {
		commandString += value
		commandString += " "
	}
	return commandString
}

func VerifyAdmin(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	var authorisedIDAdmin = "513275201375698954"
	var authorisedIDMod = "577128133975867398"

	var authorID = m.Member.Roles
	var authorised bool

	for i := range authorID {
		if authorID[i] == authorisedIDMod || authorID[i] == authorisedIDAdmin {
			authorised = true
			fmt.Println("[OK] Command authorised")
		}
	}

	return authorised
}

func VerifyTrusted(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	var authorisedIDTrusted1 = "749642547001032816"
	var authorisedIDTrusted2 = "749642583344414740"
	var authorID = m.Member.Roles
	var authorised bool

	for i := range authorID {
		if authorID[i] == authorisedIDTrusted1 || authorID[i] == authorisedIDTrusted2 {
			authorised = true
			fmt.Println("[OK] Command authorised")
		}
	}

	return authorised
}

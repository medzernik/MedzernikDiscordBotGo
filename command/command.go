// Package command Serves to regex and fix command inputs from users. Returns processed string arrays with the command and arguments parts.
package command

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"regexp"
	"strconv"
	"strings"
)

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

// ParseCommand parses a command to remove any extra fields and spaces
func ParseCommand(s string) (Command, error) {
	if !strings.HasPrefix(s, config.Cfg.ServerInfo.Prefix) && len(s) < len(config.Cfg.ServerInfo.Prefix)+1 {
		return Command{}, errors.New("")
	}

	// Remove double white spaces
	space := regexp.MustCompile(`\s+`)
	s = space.ReplaceAllString(s, " ")

	fields := strings.Fields(s)

	cmd := Command{fields[0][len(config.Cfg.ServerInfo.Prefix):], fields[1:]}

	return cmd, nil
}

// ParseMentionToString Parses the <@!userId> into userId and returns the string
func ParseMentionToString(s string) string {
	s = strings.Replace(s, "<", "", 1)
	s = strings.Replace(s, ">", "", 1)
	s = strings.Replace(s, "!", "", 1)
	s = strings.Replace(s, "@", "", 1)

	return s
}

func ParseRoleMentionToString(s string) string {
	s = strings.Replace(s, "<", "", 1)
	s = strings.Replace(s, ">", "", 1)
	s = strings.Replace(s, "!", "", 1)
	s = strings.Replace(s, "@", "", 1)
	s = strings.Replace(s, "&", "", 1)

	return s
}

func ParseStringToRoleMention(s string) string {
	var mentionID string

	mentionID = "<@&" + s + ">"

	return mentionID
}

// ParseChannelToString Parses the <#channelId> into channelId and returns the string
func ParseChannelToString(s string) string {
	s = strings.Replace(s, "<", "", 1)
	s = strings.Replace(s, ">", "", 1)
	s = strings.Replace(s, "#", "", 1)

	return s
}

// ParseStringToMentionID Parses string to output a mentionID string that Discord formats
func ParseStringToMentionID(s string) string {
	var mentionID string

	mentionID = "<@" + s + ">"

	return mentionID
}

// ParseStringToChannelID Parses string to output a channelID string that Discord formats
func ParseStringToChannelID(s string) string {
	var channelID string

	channelID = "<#" + s + ">"

	return channelID
}

// IsCommand Verifies if the command is of a command struct
func IsCommand(c *Command, name string) bool {
	return c.Command == name
}

// VerifyArguments Function verifies arguments if they are the correct format and then parses them into a struct of the command name in string and array of arguments in strings
func VerifyArguments(c *Command, args ...interface{}) error {
	if len(c.Arguments) != len(args) {
		return errors.New(c.Command + "**[ERR]** Incorrect command arguments")
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

// VerifyAdmin Function takes a bool and returns true or false based on whether the user has the admin role ID or not. Logs to stdout.
func VerifyAdmin(s *discordgo.Session, m *discordgo.MessageCreate, authorised *bool, cmd Command) bool {
	var authorID = m.Member.Roles

	//check if the user is admin, log the request if successful and then return true
	for i := range authorID {
		if authorID[i] == config.Cfg.RoleAdmin.RoleModID || authorID[i] == config.Cfg.RoleAdmin.RoleAdminID {
			*authorised = true
			fmt.Println("[OK] Command: " + cmd.Command + " authorised (Admin) by: " + m.Author.Username + " (ID: " + m.Author.ID + ")")
			break
		}
	}

	return *authorised
}

// VerifyTrusted Function takes a bool and returns true or false based on whether the user has a priviledged role (defined by admins) or not. Logs to stdout.
func VerifyTrusted(s *discordgo.Session, m *discordgo.MessageCreate, authorised *bool, cmd Command) bool {

	var authorID = m.Member.Roles

	//check if the user is trusted, log the request if successful and then return true
	for i := range authorID {
		if authorID[i] == config.Cfg.RoleTrusted.RoleTrustedID2 || authorID[i] == config.Cfg.RoleTrusted.RoleTrustedID3 {
			*authorised = true
			fmt.Println("[OK] Command: " + cmd.Command + " authorised (Trusted) by: " + m.Author.Username + " (ID: " + m.Author.ID + ")")
			break
		}
	}

	return *authorised
}

// SendTextEmbed Function will take text string of sorts and output an embed
func SendTextEmbed(s *discordgo.Session, m *discordgo.MessageCreate, status string, messageContent string, mode discordgo.EmbedType) {

	//Fixed author message.
	author := discordgo.MessageEmbedAuthor{
		URL:          "",
		Name:         "SlovakiaBot",
		IconURL:      "https://cdn.discordapp.com/avatars/837982234597916672/51236a8235b1778f5d90bce35fbcf4d6.webp?size=256",
		ProxyIconURL: "",
	}

	var color int

	//set the embed color according to the type of Status passed (OK, ERR, WARN, SYNTAX)
	switch status {
	case ":bangbang: ERROR":
		color = 15158332
	case ":warning: WARNING":
		color = 15105570
	case ":question: SYNTAX":
		color = 3447003
	case ":no_entry: AUTHENTICATION":
		color = 15105570
	case ":wrench: AUTOCORRECTING":
		color = 16776960

	default:
		color = 3066993
	}

	//MessageEmbed info
	//Thinking of adding timestamp time.Now().Format(time.RFC3339)
	embed := discordgo.MessageEmbed{
		URL:         "",
		Type:        mode,
		Title:       status,
		Description: messageContent,
		Timestamp:   "",
		Color:       color,
		Footer:      nil,
		Image:       nil,
		Thumbnail:   nil,
		Video:       nil,
		Provider:    nil,
		Author:      &author,
		Fields:      nil,
	}

	//Send a message as an embed.
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "**[TESTERR]** Error: "+err.Error())
		return
	}
}
func SendTextEmbedCommand(s *discordgo.Session, m string, status string, messageContent string, mode discordgo.EmbedType) {

	//Fixed author message.
	author := discordgo.MessageEmbedAuthor{
		URL:          "",
		Name:         "SlovakiaBot",
		IconURL:      "https://cdn.discordapp.com/avatars/837982234597916672/51236a8235b1778f5d90bce35fbcf4d6.webp?size=256",
		ProxyIconURL: "",
	}

	var color int

	//set the embed color according to the type of Status passed (OK, ERR, WARN, SYNTAX)
	switch status {
	case ":bangbang: ERROR":
		color = 15158332
	case ":warning: WARNING":
		color = 15105570
	case ":question: SYNTAX":
		color = 3447003
	case ":no_entry: AUTHENTICATION":
		color = 15105570
	case ":wrench: AUTOCORRECTING":
		color = 16776960

	default:
		color = 3066993
	}

	var embedArray []discordgo.MessageEmbed
	//MessageEmbed info
	//Thinking of adding timestamp time.Now().Format(time.RFC3339)
	embed := discordgo.MessageEmbed{
		URL:         "",
		Type:        mode,
		Title:       status,
		Description: messageContent,
		Timestamp:   "",
		Color:       color,
		Footer:      nil,
		Image:       nil,
		Thumbnail:   nil,
		Video:       nil,
		Provider:    nil,
		Author:      &author,
		Fields:      nil,
	}

	embedArray = append(embedArray, embed)

	//Send a message as an embed.
	_, err := s.ChannelMessageSendEmbed(m, &embed)
	if err != nil {
		s.ChannelMessageSend(m, "**[TESTERR]** Error: "+err.Error())
		return
	}
}

func VerifyAdminCMD(s *discordgo.Session, m string, authorised *bool, cmd *discordgo.InteractionCreate) bool {
	var authorID = cmd.Member.Roles

	//check if the user is admin, log the request if successful and then return true
	for i := range authorID {
		if authorID[i] == config.Cfg.RoleAdmin.RoleModID || authorID[i] == config.Cfg.RoleAdmin.RoleAdminID {
			*authorised = true
			fmt.Println("[OK] Command" + " authorised (Admin) by: " + cmd.Member.Nick + " (ID: " + cmd.Member.User.ID + ")")
			break
		}
	}

	return *authorised
}

// VerifyTrusted Function takes a bool and returns true or false based on whether the user has a priviledged role (defined by admins) or not. Logs to stdout.
func VerifyTrustedCMD(s *discordgo.Session, m string, authorised *bool, cmd *discordgo.InteractionCreate) bool {

	var authorID = cmd.Member.Roles

	//check if the user is trusted, log the request if successful and then return true
	for i := range authorID {
		if authorID[i] == config.Cfg.RoleTrusted.RoleTrustedID2 || authorID[i] == config.Cfg.RoleTrusted.RoleTrustedID3 {
			*authorised = true
			fmt.Println("[OK] Command" + " authorised (Trusted) by: " + cmd.Member.Nick + " (ID: " + cmd.Member.User.ID + ")")
			break
		}
	}

	return *authorised
}

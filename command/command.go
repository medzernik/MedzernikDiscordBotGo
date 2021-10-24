// Package command Serves to regex and fix command inputs from users. Returns processed string arrays with the command and arguments parts.
package command

import (
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"strings"
	"time"
)

// ParseMentionToString Parses the <@!userId> into userId and returns the string
func ParseMentionToString(s string) string {
	s = strings.Replace(s, "<", "", 1)
	s = strings.Replace(s, ">", "", 1)
	s = strings.Replace(s, "!", "", 1)
	s = strings.Replace(s, "@", "", 1)

	return s
}

// ParseRoleMentionToString Parses the <@&roleId> into roleId and returns the string
func ParseRoleMentionToString(s string) string {
	s = strings.Replace(s, "<", "", 1)
	s = strings.Replace(s, ">", "", 1)
	s = strings.Replace(s, "!", "", 1)
	s = strings.Replace(s, "@", "", 1)
	s = strings.Replace(s, "&", "", 1)

	return s
}

// ParseStringToRoleMention Parses the string Id and returns <@&roleId>
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

// SendTextEmbedCommand Custom embed parser. Sends a parsed embed.
func SendTextEmbedCommand(s *discordgo.Session, m string, status string, messageContent string, mode discordgo.EmbedType) {
	//Fixed author message.
	author := discordgo.MessageEmbedAuthor{
		URL:          "",
		Name:         config.Cfg.ServerInfo.BotName,
		IconURL:      config.Cfg.ServerInfo.BotLogo,
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
		//Set a default color as well
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
		Timestamp:   time.Now().Format(time.RFC3339),
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

// MemberHasPermission Check whether the member has permissions of a given number (in integer)
func MemberHasPermission(s *discordgo.Session, guildID string, userID string, permission int64) (bool, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		if member, err = s.GuildMember(guildID, userID); err != nil {
			return false, err
		}
	}

	// Iterate through the role IDs stored in member.Roles
	// to check permissions
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			return false, err
		}
		if role.Permissions&permission != 0 {
			return true, nil
		}
	}

	return false, nil
}

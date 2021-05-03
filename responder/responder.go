package responder

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func RegisterPlugin(s *discordgo.Session) {
	s.AddHandler(messageCreated)
	s.AddHandler(reactionAdded)
	s.AddHandler(ready)

}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "Zasielkovna" {
		s.ChannelMessageSend(m.ChannelID, "OVER 200% <a:medzernikShake:814055147583438848>")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func reactionAdded(s *discordgo.Session, mr *discordgo.MessageReactionAdd) {
	if strings.Contains(strings.ToLower(mr.Emoji.Name), "kekw") {
		s.MessageReactionAdd(mr.ChannelID, mr.MessageID, mr.Emoji.APIName())
	}
}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the status.
	s.UpdateGameStatus(5, "test")
}

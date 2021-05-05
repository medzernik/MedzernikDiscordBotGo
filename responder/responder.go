package responder

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
	"time"
)

func RegisterPlugin(s *discordgo.Session) {
	s.AddHandler(messageCreated)
	s.AddHandler(reactionAdded)
	s.AddHandler(ready)
	s.AddHandler(SnowflakeTimestamp)

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

	if m.Content == "-rayman" {
		rayman_time_raw, _ := SnowflakeTimestamp("242318670079066112")

		rayman_time := time.Now().Sub(rayman_time_raw)
		rayman_time_days := rayman_time.Hours() / 24
		rayman_time_days = rayman_time_days / 24

		//rayman_time_string := strconv.FormatFloat(rayman_time, 'f', 5, 64)
		rayman_time_days_string := rayman_time.Hours() / 24
		fmt.Println(rayman_time_days_string)
		rayman_time_days_string_pure := strconv.FormatFloat(rayman_time_days_string, 'f', 0, 64)

		rayman_time_string := rayman_time.String()
		rayman_time_string = strings.ReplaceAll(rayman_time_string, "h", " Hodín\n ")
		rayman_time_string = strings.ReplaceAll(rayman_time_string, "m", " Minút\n ")
		rayman_time_string = strings.ReplaceAll(rayman_time_string, "s", " Sekúnd ")

		fmt.Println("log -rayman: ", rayman_time_string)

		s.ChannelMessageSend(m.ChannelID, "Rayman je tu s nami už:\n "+rayman_time_days_string_pure+" (celkovo dní), rozpis:\n----------\n "+rayman_time_string+"<:peepoLove:687313976043765810>")
	}

	if m.Content == "-check users" {
		members, _ := s.GuildMembers("513274646406365184", "0", 1000)

		for itera, _ := range members {
			user_time_join, _ := members[itera].JoinedAt.Parse()

			timevar := user_time_join.Sub(time.Now()).Hours()
			fmt.Println(timevar)

			if timevar > -24 {
				println("THIS USER IS TOO YOUNG")

				s.ChannelMessageSend(m.ChannelID, "This user is too young (less than 24h join age): "+members[itera].User.Username)
			}

			//fmt.Printf("difference = %v\n", difference)

			//println(usertime)
		}

		fmt.Println(len(members))

		//s.ChannelMessageSend(m.ChannelID, user_warn)

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
	s.UpdateGameStatus(0, "Welcome to Slovakia")

}

func SnowflakeTimestamp(ID string) (t time.Time, err error) {
	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	timestamp := (i >> 22) + 1420070400000
	t = time.Unix(0, timestamp*1000000)
	fmt.Println(t)
	return
}

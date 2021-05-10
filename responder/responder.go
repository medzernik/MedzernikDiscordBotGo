package responder

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
)

func RegisterPlugin(s *discordgo.Session) {
	s.AddHandler(messageCreated)
	//s.AddHandler(reactionAdded)
	s.AddHandler(ready)

}
//SETTINGS
//guildID. Change this to represent your server.
var guildidnumber = "513274646406365184"


//This is the main logic and command file for now
//TODO: implement a system of an internal user database

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itsel
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	cmd, err := command.ParseCommand(m.Content)

	if err != nil {
		println(err.Error())
		return
	}

	// If the message is "ping" reply with "Pong!"
	if command.IsCommand(&cmd, "Zasielkovna") {
		err := command.VerifyArguments(&cmd)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		s.ChannelMessageSend(m.ChannelID, "OVER 200% <a:medzernikShake:814055147583438848>")

	}

	//a personal reward for our founder of the server that tracks his time on the guilds
	if command.IsCommand(&cmd, "age") {
		err := command.VerifyArguments(&cmd, command.RegexArg{`^<@!(\d+)>$`, 1})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		userId := cmd.Arguments[0]
		//Every time a command is run, get a list of all users. This serves the purpose to then print the name of the corresponding user.
		//TODO: Make this list a (ideally cached) variable that at least is shared and not run every time a command is run.
		membersCached := GetMemberListFromGuild(s, guildidnumber)

		var userName string

		for itera := range membersCached {
			if membersCached[itera].User.ID == userId {
				userName = membersCached[itera].User.Username
				fmt.Println(userName)
			}
		}

		fmt.Println(userId)

		userTimeRaw, err := SnowflakeTimestamp(userId)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Zlé ID slováka")

			return
		}

		if userId < "0" {
			return
		}

		userTime := time.Now().Sub(userTimeRaw)
		userTimeDays := userTime.Hours() / 24
		userTimeDays = userTimeDays / 24

		userTimeDaysString := userTime.Hours() / 24
		fmt.Println(userTimeDaysString)
		userTimeDaysStringPure := strconv.FormatFloat(userTimeDaysString, 'f', 0, 64)

		userTimeString := userTime.String()
		userTimeString = strings.ReplaceAll(userTimeString, "h", " Hodín\n ")
		userTimeString = strings.ReplaceAll(userTimeString, "m", " Minút\n ")
		userTimeString = strings.ReplaceAll(userTimeString, "s", " Sekúnd ")

		fmt.Println("log -rayman: ", userTimeString)

		s.ChannelMessageSend(m.ChannelID, "**"+userName+"**"+" je tu s nami už:\n "+userTimeDaysStringPure+" (celkovo dní), rozpis:\n----------\n "+userTimeString+"<:peepoLove:687313976043765810>")

	}


	//right now this command checks for any 1000 users on the guild that have a join time less than 24hours, then prints the names one by one.
	//TODO: check if the users can be >1000
	//TODO: implement a raid protection checker that checks every 1 hour for accounts <2 hours of age and if finds more than 5 -> alert the admins
	//TODO [DONE]: change the output message to be a single message in a single output to protect from spam. Change the information.
	if command.IsCommand(&cmd, "check-users") {
		err := command.VerifyArguments(&cmd)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}
		//variable definitons
		//members_cached, _ := s.GuildMembers("513274646406365184", "0", 1000)
		membersCached := GetMemberListFromGuild(s, guildidnumber)
		var tempMsg string
		timeToCheckUsers := 24.0 * -1.0

		//iterate over the members_cached array. Maximum limit is 1000.
		for itera := range membersCached {
			userTimeJoin, _ := membersCached[itera].JoinedAt.Parse()
			timevar := userTimeJoin.Sub(time.Now()).Hours()

			fmt.Println(timevar)

			if timevar > timeToCheckUsers {
				println("THIS USER IS TOO YOUNG")

				tempMsg += "This user is too young (less than 24h join age): " + membersCached[itera].User.Username + "\n"
			}
		}
		//print out the amount of members_cached (max is currently 1000)
		fmt.Println(len(membersCached))
		s.ChannelMessageSend(m.ChannelID, tempMsg)
	}

}

/*
//this function adds a +1 to a specific emoji reaction to an already added one by a use
//TODO: make it a bit more modular and expand the amount of reactions. Ideally a variable level system
func reactionAdded(s *discordgo.Session, mr *discordgo.MessageReactionAdd) {
	if strings.ToUpper(mess) == , "kekw") {

		s.MessageReactionAdd(mr.ChannelID, mr.MessageID, mr.Emoji.APIName())
	}
	if strings.Contains(strings.ToLower(mr.Emoji.Name), "okayChamp") {
		s.MessageReactionAdd(mr.ChannelID, mr.MessageID, mr.Emoji.APIName())
	}

}
*/



// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the status.
	s.UpdateGameStatus(0, "Gde mozog")

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

func GetMemberListFromGuild(s *discordgo.Session, guildID string) []*discordgo.Member {
	membersList, err := s.GuildMembers(guildID, "0", 1000)
	if err != nil {
		fmt.Println("Error getting the memberlist (probably invalid guildID): " + guildID)
	}

	return membersList

}

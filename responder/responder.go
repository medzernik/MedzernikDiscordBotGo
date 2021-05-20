package responder

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"strings"

	//_ "modernc.org/sqlite"
	"strconv"
	"time"
)

func RegisterPlugin(s *discordgo.Session) {
	s.AddHandler(messageCreated)
	//s.AddHandler(reactionAdded)
	s.AddHandler(ready)

}

//SETTINGS
//guildID. Change this to represent your server. ID of channel and server, String data type
var guildidnumber = "513274646406365184"
var adminchannel = "837987736416813076"
var logchannel = "513280604507340804"

//This is the main logic and command file for now
//TODO: implement a system of an internal user database (redis?)

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

	//Time since joined command
	if command.IsCommand(&cmd, "age") {
		err := command.VerifyArguments(&cmd, command.RegexArg{Expression: `^<@!(\d+)>$`, CaptureGroup: 1})
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		userId := cmd.Arguments[0]
		//Every time a command is run, get a list of all users. This serves the purpose to then print the name of the corresponding user.
		//TODO: cache it in redis
		membersCached := GetMemberListFromGuild(s, guildidnumber)

		var userName string

		for itera := range membersCached {
			if membersCached[itera].User.ID == userId {
				userName = membersCached[itera].User.Username
				fmt.Println(userName)
			} else if membersCached[itera].User.ID != userId && membersCached[itera].User.ID == "" {
				s.ChannelMessageSend(m.ChannelID, "Zlé ID slováka")

			}
		}

		fmt.Println(userId)

		userTimeRaw, err := SnowflakeTimestamp(userId)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Zlé ID slováka")
			return
		}

		userTime := time.Now().Sub(userTimeRaw)

		dny := int64(userTime.Hours() / 24)
		hodiny := int64(userTime.Hours()) - dny*24
		minuty := int64(userTime.Minutes()) - int64(userTime.Hours())*60
		sekundy := int64(userTime.Seconds()) - int64(userTime.Minutes())*60

		dnyString := strconv.FormatInt(dny, 10)
		hodinyString := strconv.FormatInt(hodiny, 10)
		minutyString := strconv.FormatInt(minuty, 10)
		sekundyString := strconv.FormatInt(sekundy, 10)

		s.ChannelMessageSend(m.ChannelID, "**"+userName+"**"+" je tu s nami už:\n"+dnyString+" dni\n"+hodinyString+" hodin\n"+minutyString+" minut\n"+sekundyString+" sekund"+"<:peepoLove:687313976043765810>")

	}
	if command.IsCommand(&cmd, "mute") {
		var muteUser string
		var timeMute uint64
		var roleMuteID string = "684159104901709825"
		var priviledgedRole string = "513275201375698954"
		membersCached := GetMemberListFromGuild(s, guildidnumber)

		if len(cmd.Arguments) < 1 {
			_, err := s.ChannelMessageSend(m.ChannelID, "Insufficient arguments provded (help: .mute **@user time**)")
			if err != nil {
				return
			}
			return
		}
		muteUserRaw := cmd.Arguments[0]
		muteUser = command.ParseMentionToString(cmd.Arguments[0])

		if len(cmd.Arguments) > 1 {
			timeMute, err = strconv.ParseUint(cmd.Arguments[1], 10, 64)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, "**Error parsing the time** (help: second argument should be number in minutes, max 64-bit sized.) \nExample: .mute user **30**"+"\nYou can also use the command without time specified"+"\n"+err.Error())
				return
			}
		}

		for itera := range membersCached {
			if membersCached[itera].User.ID == muteUser {
				for i := range membersCached[itera].Roles {
					if membersCached[itera].Roles[i] == roleMuteID {
						s.ChannelMessageSend(m.ChannelID, "**Duplicate request:** User "+muteUser+" already has a role Muted "+roleMuteID+" added")
						return
					} else if muteUser == "837982234597916672" {
						s.ChannelMessageSend(m.ChannelID, "***You fucking donkey. WHO DO YOU THINK YOU ARE MUTING YOU FUCK*** <:pocem:683715405407059968> <:pocem:683715405407059968> <:pocem:683715405407059968> <:pocem:683715405407059968>")
						return
					} else if membersCached[itera].Roles[i] == priviledgedRole && membersCached[itera].User.ID == muteUser {
						s.ChannelMessageSend(m.ChannelID, "**Refusing to mute user** "+muteUserRaw+" **who has a privileged role** "+priviledgedRole+"\n The role is set as priviledged and can't be assigned role id: "+roleMuteID)
						return

					}

				}

			}
		}
		for itera := range membersCached {
			for i := range membersCached[itera].Roles {
				if membersCached[itera].User.ID == m.Author.ID && membersCached[itera].Roles[i] == priviledgedRole {
					var err error
					err = s.GuildMemberMute(guildidnumber, muteUser, true)
					if err != nil {
					}
					err = s.GuildMemberRoleAdd(guildidnumber, muteUser, roleMuteID)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "Error adding muted role to user "+muteUserRaw+"\n"+err.Error())
						return
					}
					s.ChannelMessageSend(m.ChannelID, "User muted successfully.")
					s.ChannelMessageSend(logchannel, "User "+muteUserRaw+" muted successfully - assigned role "+roleMuteID)
					return
				} else {
					s.ChannelMessageSend(m.ChannelID, "Insufficient permissions")
					fmt.Println(timeMute)
					return
				}

			}
		}
	}

	//right now this command checks for any 1000 users on the guild that have a join time less than 24hours, then prints the names one by one.
	//TODO: check if the users can be >1000
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

			if timevar > timeToCheckUsers {
				tempMsg += "This user is too young (less than 24h join age): " + membersCached[itera].User.Username + "\n"
			}
			fmt.Println("Found ", strings.Contains(tempMsg, "\n"), " very young users...")
		}
		//print out the amount of members_cached (max is currently 1000)
		s.ChannelMessageSend(m.ChannelID, tempMsg)

	}

	if command.IsCommand(&cmd, "plan") {
		if len(cmd.Arguments) < 3 {
			s.ChannelMessageSend(m.ChannelID, "Insufficient arguments. Provided "+strconv.FormatInt(int64(len(cmd.Arguments)), 10)+" , Expected at least 3")
			return
		}
		GamePlanner(&cmd)
	}

}

func GamePlanner(c *command.Command) {

	//db.Ping()

	type PlannedGame struct {
		ID          uint `gorm:"primaryKey"`
		Name        string
		TimePlanned string
		Mentions    string
		CreatedAt   time.Time
		UpdatedAt   time.Time
	}
	return
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	//set the status
	err := s.UpdateGameStatus(0, "Gde mozog")
	if err != nil {
		return
	}
	//run the raid checker function
	go CheckRegularSpamAttack(s)
	// github.com/mattn/go-sqlite3

}

//var db, err = sql.Open("sqlite", dsnURI)

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

func CheckRegularSpamAttack(s *discordgo.Session) {
	//variable definitons
	var membersCached = GetMemberListFromGuild(s, guildidnumber)
	var tempMsg string
	var spamcounter int64
	var checkinterval time.Duration = 90
	var timeToCheckUsers = 10 * -1.0

	for {
		//iterate over the members_cached array. Maximum limit is 1000.
		for itera := range membersCached {
			userTimeJoin, _ := membersCached[itera].JoinedAt.Parse()
			timevar := userTimeJoin.Sub(time.Now()).Minutes()

			if timevar > timeToCheckUsers {
				println("Checking " + membersCached[itera].User.Username + " for possible raid attack")
				tempMsg += "RAID PROTECTION ALERT!: User" + membersCached[itera].User.Username + "join age: " + strconv.FormatFloat(timeToCheckUsers, 'f', 0, 64) + "\n"
				spamcounter += 1
			}

		}
		if spamcounter > 4 {
			s.ChannelMessageSend(adminchannel, "WARN: Possible RAID ATTACK detected!!! (<@&513275201375698954>) ("+strconv.FormatInt(spamcounter, 10)+" users joined in the last "+strconv.FormatFloat(timeToCheckUsers, 'f', 0, 64)+" hours)")
		}
		spamcounter = 0
		time.Sleep(checkinterval * time.Second)
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

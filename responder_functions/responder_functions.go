// Package responder_functions contains all the logic and basic config commands for the responder commands.
package responder_functions

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"github.com/medzernik/SlovakiaDiscordBotGo/database"
	"strconv"
	"strings"
	"time"
)

const Version string = "0.5.1"
const VersionFeatureName string = "The Slowmode Update"

type CommandStatus struct {
	OK      string
	ERR     string
	SYNTAX  string
	WARN    string
	AUTH    string
	AUTOFIX string
}

// CommandStatusBot is a variable to pass to the messageEmbed to make an emoji
var CommandStatusBot CommandStatus = CommandStatus{
	OK:      "",
	ERR:     ":bangbang: ERROR",
	SYNTAX:  ":question: SYNTAX",
	WARN:    ":warning: WARNING",
	AUTH:    ":no_entry: AUTHENTICATION",
	AUTOFIX: ":wrench: AUTOCORRECTING",
}

// GamePlanInsert Inserts the game into the database
func GamePlanInsert(c *command.Command, s **discordgo.Session, m **discordgo.MessageCreate) {
	//open database and then close it (defer)
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer func(sqliteDatabase *sql.DB) {
		err := sqliteDatabase.Close()
		if err != nil {

		}
	}(sqliteDatabase)

	//transform to timestamp
	splitTimeArgument := strings.Split(c.Arguments[0], ":")

	//TODO: Check the capacity if it's sufficient, otherwise the program is panicking every time...
	if cap(splitTimeArgument) < 1 {
		command.SendTextEmbed(*s, *m, CommandStatusBot.ERR, "Error parsing time", discordgo.EmbedTypeRich)

		//(*s).ChannelMessageSend((*m).ChannelID, "**[ERR]** Error parsing time")
		return
	}

	//Put hours into timeHours
	timeHour, err := strconv.Atoi(splitTimeArgument[0])
	if err != nil {
		command.SendTextEmbed(*s, *m, CommandStatusBot.ERR, "Error converting hours", discordgo.EmbedTypeRich)
		//(*s).ChannelMessageSend((*m).ChannelID, "**[ERR]** Error converting hours")
		//fmt.Printf("%s", err)
		return
	}
	//put minutes into timeMinute
	timeMinute, err := strconv.Atoi(splitTimeArgument[1])
	if err != nil {
		command.SendTextEmbed(*s, *m, CommandStatusBot.ERR, "Error converting minutes", discordgo.EmbedTypeRich)
		//(*s).ChannelMessageSend((*m).ChannelID, "**[ERR]** Error converting minutes")
		//fmt.Printf("%s", err)
		return
	}
	//get current date and replace hours and minutes with user variables
	gameTimestamp := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), timeHour, timeMinute, time.Now().Second(), 0, time.Now().Location())
	gameTimestampInt := gameTimestamp.Unix()

	fmt.Println(gameTimestampInt)

	//export to database
	database.InsertGame(sqliteDatabase, gameTimestampInt, c.Arguments[1], c.Arguments[2])

	var plannedGames string
	database.DisplayGamePlanned(sqliteDatabase, &plannedGames)

	command.SendTextEmbed(*s, *m, CommandStatusBot.OK+"PLANNED A GAME", plannedGames, discordgo.EmbedTypeRich)
	return
}

// SnowflakeTimestamp Function to check the user's join date
func SnowflakeTimestamp(ID string) (t time.Time, err error) {
	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	timestamp := (i >> 22) + 1420070400000
	t = time.Unix(0, timestamp*1000000)
	return
}

// GetGuildInfo gets info of a guild (preview) and returns the struct
func GetGuildInfo(s *discordgo.Session, guildIDT string) *discordgo.GuildPreview {
	guildInfo, err := s.GuildPreview(guildIDT)
	if err != nil {
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error getting preview info about guild: "+guildIDT)
	}
	return guildInfo
}

// GetMemberListFromGuild Gets the member info and tries to save it in a local (cached sort of) array to access later.
func GetMemberListFromGuild(s *discordgo.Session, guildID string) []*discordgo.Member {
	//TODO: This only works as a preview.. wtf
	guildInfoTemp := GetGuildInfo(s, config.Cfg.ServerInfo.GuildIDNumber)

	membersList, err := s.GuildMembers(guildID, "0", guildInfoTemp.ApproximateMemberCount)
	if err != nil {
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error getting information about users with the guildID. InvalidID or >1000 members (bug)")
		fmt.Println("ERROR: ", err)
		return membersList
	}

	//TODO: Fix this? How? Bug in the library?
	/*
		//membersListAdditional, _ := s.GuildMembers(guildID, "1000", 1000)


		for j := range membersListAdditional {
			membersList = append(membersList, membersListAdditional[j])
		}

	*/

	return membersList
}

// CheckRegularSpamAttack Checks the server for spam attacks
func CheckRegularSpamAttack(s *discordgo.Session) {
	//variable definitons
	var membersCached = GetMemberListFromGuild(s, config.Cfg.ServerInfo.GuildIDNumber)
	var tempMsg string
	var spamCounter int64
	var checkInterval time.Duration = 90
	var timeToCheckUsers = 10 * -1.0

	for {
		//iterate over the members_cached array. Maximum limit is 1000.
		for i := range membersCached {
			userTimeJoin, _ := membersCached[i].JoinedAt.Parse()
			timeVar := userTimeJoin.Sub(time.Now()).Minutes()

			if timeVar > timeToCheckUsers {
				tempMsg += "**[ALERT]** RAID PROTECTION ALERT!: User" + membersCached[i].User.Username + "join age: " + strconv.FormatFloat(timeToCheckUsers, 'f', 0, 64) + "\n"
				spamCounter += 1
			}

		}
		if spamCounter > 4 {
			s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[WARN]** Possible RAID ATTACK detected!!! (<@&513275201375698954>) ("+command.ParseStringToMentionID(config.Cfg.RoleAdmin.RoleAdminID+strconv.FormatInt(spamCounter, 10)+" users joined in the last "+strconv.FormatFloat(timeToCheckUsers, 'f', 0, 64)+" hours)"))
		}
		spamCounter = 0
		time.Sleep(checkInterval * time.Second)
	}

}

func Fox(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "<a:medzernikShake:814055147583438848>")
	return
}

// TimedChannelUnlock automatically locks and unlocks a trusted channel
func TimedChannelUnlock(s *discordgo.Session) {
	if config.Cfg.AutoLocker.Enabled == false {
		return
	}

	var checkInterval time.Duration = 60

	fmt.Println("[INIT OK] Channel unlock system module initialized")

	for {
		if time.Now().Weekday() == config.Cfg.AutoLocker.TimeDayUnlock && time.Now().Hour() == config.Cfg.AutoLocker.TimeHourUnlock && time.Now().Minute() == config.Cfg.AutoLocker.TimeMinuteUnlock {
			//Unlock the channel
			//TargetType 0 = roleID, 1 = memberID
			err1 := s.ChannelPermissionSet(config.Cfg.RoleTrusted.ChannelTrustedID, config.Cfg.RoleTrusted.RoleTrustedID1, 0, 2251673408, 0)
			if err1 != nil {
				s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error changing the permissions for role "+"<@"+config.Cfg.RoleTrusted.RoleTrustedID1+">")
			}
			err2 := s.ChannelPermissionSet(config.Cfg.RoleTrusted.ChannelTrustedID, config.Cfg.RoleTrusted.RoleTrustedID2, 0, 2251673408, 0)
			if err2 != nil {
				s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error changing the permissions for role "+"<@"+config.Cfg.RoleTrusted.RoleTrustedID2+">")
			}
			err3 := s.ChannelPermissionSet(config.Cfg.RoleTrusted.ChannelTrustedID, config.Cfg.RoleTrusted.RoleTrustedID3, 0, 2251673408, 0)
			if err3 != nil {
				s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error changing the permissions for role "+"<@"+config.Cfg.RoleTrusted.RoleTrustedID3+">")
			}
			err4 := s.ChannelPermissionSet(config.Cfg.RoleTrusted.ChannelTrustedID, config.Cfg.RoleTrusted.RoleTrustedID4, 0, 2251673408, 0)
			if err4 != nil {
				s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error changing the permissions for role "+"<@"+config.Cfg.RoleTrusted.RoleTrustedID4+">")
			}

			fmt.Println("[OK] Opened the channel " + config.Cfg.RoleTrusted.ChannelTrustedID)
		} else if time.Now().Weekday() == config.Cfg.AutoLocker.TimeDayLock && time.Now().Hour() == config.Cfg.AutoLocker.TimeHourLock && time.Now().Minute() == config.Cfg.AutoLocker.TimeMinuteLock {
			//Lock the channel
			//TargetType 0 = roleID, 1 = memberID
			err1 := s.ChannelPermissionSet(config.Cfg.RoleTrusted.ChannelTrustedID, config.Cfg.RoleTrusted.RoleTrustedID1, 0, 0, 2251673408)
			if err1 != nil {
				s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error changing the permissions for role "+"<@"+config.Cfg.RoleTrusted.RoleTrustedID1+">")
			}
			err2 := s.ChannelPermissionSet(config.Cfg.RoleTrusted.ChannelTrustedID, config.Cfg.RoleTrusted.RoleTrustedID2, 0, 0, 2251673408)
			if err2 != nil {
				s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error changing the permissions for role "+"<@"+config.Cfg.RoleTrusted.RoleTrustedID2+">")
			}
			err3 := s.ChannelPermissionSet(config.Cfg.RoleTrusted.ChannelTrustedID, config.Cfg.RoleTrusted.RoleTrustedID3, 0, 0, 2251673408)
			if err3 != nil {
				s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error changing the permissions for role "+"<@"+config.Cfg.RoleTrusted.RoleTrustedID3+">")
			}
			err4 := s.ChannelPermissionSet(config.Cfg.RoleTrusted.ChannelTrustedID, config.Cfg.RoleTrusted.RoleTrustedID4, 0, 0, 2251673408)
			if err4 != nil {
				s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error changing the permissions for role "+"<@"+config.Cfg.RoleTrusted.RoleTrustedID4+">")
			}
			fmt.Println("[OK] Closed the channel " + config.Cfg.RoleTrusted.ChannelTrustedID)
		}

		time.Sleep(checkInterval * time.Second)
	}

}

// Members returns a value of members on the server currently to another function.
func Members(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) uint64 {
	if len(cmd.Arguments) > 0 {
		command.SendTextEmbed(s, m, CommandStatusBot.SYNTAX, "Usage: **.count**\n Automatically discarding arguments...", discordgo.EmbedTypeRich)
	}

	memberList := GetMemberListFromGuild(s, config.Cfg.ServerInfo.GuildIDNumber)
	memberListLength := uint64(len(memberList))

	return memberListLength
}

// PruneCount returns a value to another function of how many members to prune
func PruneCount(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) uint32 {

	if len(cmd.Arguments) < 1 {
		command.SendTextEmbed(s, m, CommandStatusBot.SYNTAX, "Usage **.prunecount days**", discordgo.EmbedTypeRich)
		return 0
	}

	pruneDaysString := cmd.Arguments[0]
	pruneDaysInt, err1 := strconv.ParseInt(pruneDaysString, 10, 64)
	if err1 != nil {
		command.SendTextEmbed(s, m, CommandStatusBot.ERR, "Error parsing the argument as uint32 days number", discordgo.EmbedTypeRich)
		return 0
	}

	if pruneDaysInt < 7 {
		pruneDaysInt = 0
		command.SendTextEmbed(s, m, CommandStatusBot.WARN, "Command is limited to range 7-30 for safety reasons", discordgo.EmbedTypeRich)
		return 0
	}

	pruneDaysCount, err2 := s.GuildPruneCount(config.Cfg.ServerInfo.GuildIDNumber, uint32(pruneDaysInt))
	if err2 != nil {
		return 0
	}
	return pruneDaysCount

}

// MassKick mass kicks user IDs
func MassKick(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {

}

// MassBan mass bans user IDs
func MassBan(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {

}

// ConfigurationReload reloads the config file
func ConfigurationReload(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	var authorisedAdmin bool = false
	authorisedAdmin = command.VerifyAdmin(s, m, &authorisedAdmin, cmd)
	if authorisedAdmin == true {
		config.LoadConfig()
		command.SendTextEmbed(s, m, CommandStatusBot.OK+"CONFIG LOADED", "Loaded the new config", discordgo.EmbedTypeRich)
	} else {
		command.SendTextEmbed(s, m, CommandStatusBot.AUTH, "Insufficient permissions", discordgo.EmbedTypeRich)
	}
}

func FoxTest(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.ChannelMessageSend(i.ChannelID, "<a:medzernikShake:814055147583438848>")
	return
}

// FindUserVoiceState Finds the user in a voice channel
func FindUserVoiceState(session *discordgo.Session, userid string) (*discordgo.VoiceState, error) {
	for _, guild := range session.State.Guilds {
		for _, vs := range guild.VoiceStates {
			if vs.UserID == userid {
				return vs, nil
			}
		}
	}
	return nil, errors.New("Could not find user's voice state")
}

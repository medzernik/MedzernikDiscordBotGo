// Package responder_functions contains all the logic and basic config commands for the responder commands.
package responder_functions

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"strconv"

	"time"
)

const Version string = "0.8.9"
const VersionFeatureName string = "The Logging Update (WIP)"

const TargetTypeRoleID discordgo.PermissionOverwriteType = 0
const TargetTypeMemberID discordgo.PermissionOverwriteType = 1

/*
// GamePlanInsert Inserts the game into the database
func GamePlanInsert(c *command.Command, s **discordgo.Session, m **discordgo.MessageCreate) {
	//open database and then close it (defer)
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer func(sqliteDatabase *sql.DB) {
		err := sqliteDatabase.Close()
		if err != nil {

		}
	}(sqliteDatabase)



	//get current date and replace hours and minutes with user variables
	gameTimestamp := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), timeHour, timeMinute, time.Now().Second(), 0, time.Now().Location())
	gameTimestampInt := gameTimestamp.Unix()

	fmt.Println(gameTimestampInt)

	//export to database
	database.InsertGame(sqliteDatabase, gameTimestampInt, c.Arguments[1], c.Arguments[2])

	var plannedGames string
	database.DisplayGamePlanned(sqliteDatabase, &plannedGames)

	command.SendTextEmbedCommand(*s, m, CommandStatusBot.OK+"PLANNED A GAME", plannedGames, discordgo.EmbedTypeRich)
	return
}

*/

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

/*
// GetGuildInfo gets info of a guild (preview) and returns the struct
func GetGuildInfo(s *discordgo.Session, guildID string) *discordgo.GuildPreview {
	guildInfo, err := s.GuildPreview(guildID)
	if err != nil {
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error getting preview info about guild: "+guildID)
	}
	return guildInfo
}

*/

func UnlockTrustedChannel(s *discordgo.Session, perms int64, target discordgo.PermissionOverwriteType) {
	var roleArrayToUnlock []string
	roleArrayToUnlock = append(roleArrayToUnlock, config.Cfg.RoleTrusted.RoleTrustedID1)
	roleArrayToUnlock = append(roleArrayToUnlock, config.Cfg.RoleTrusted.RoleTrustedID2)
	roleArrayToUnlock = append(roleArrayToUnlock, config.Cfg.RoleTrusted.RoleTrustedID3)
	roleArrayToUnlock = append(roleArrayToUnlock, config.Cfg.RoleTrusted.RoleTrustedID4)

	for i := 0; i < len(roleArrayToUnlock); i++ {
		//Unlock the channel
		//TargetType 0 = roleID, 1 = memberID
		err1 := s.ChannelPermissionSet(config.Cfg.RoleTrusted.ChannelTrustedID, roleArrayToUnlock[i], target, perms, 0)
		if err1 != nil {
			s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error changing the permissions for role "+command.ParseStringToRoleMention(roleArrayToUnlock[i]))
		}
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Unlocked the channel for "+command.ParseStringToRoleMention(roleArrayToUnlock[i]))
	}
	return
}
func LockTrustedChannel(s *discordgo.Session, perms int64, target discordgo.PermissionOverwriteType) {
	var roleArrayToUnlock []string
	roleArrayToUnlock = append(roleArrayToUnlock, config.Cfg.RoleTrusted.RoleTrustedID1)
	roleArrayToUnlock = append(roleArrayToUnlock, config.Cfg.RoleTrusted.RoleTrustedID2)
	roleArrayToUnlock = append(roleArrayToUnlock, config.Cfg.RoleTrusted.RoleTrustedID3)
	roleArrayToUnlock = append(roleArrayToUnlock, config.Cfg.RoleTrusted.RoleTrustedID4)

	for i := 0; i < len(roleArrayToUnlock); i++ {
		//Unlock the channel
		//TargetType 0 = roleID, 1 = memberID
		err1 := s.ChannelPermissionSet(config.Cfg.RoleTrusted.ChannelTrustedID, roleArrayToUnlock[i], target, 0, perms)
		if err1 != nil {
			s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error changing the permissions for role "+command.ParseStringToRoleMention(roleArrayToUnlock[i]))
		}
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Locked the channel for "+command.ParseStringToRoleMention(roleArrayToUnlock[i]))
	}
	return
}

var LockChannelToday bool = false

// TimedChannelUnlock automatically locks and unlocks a trusted channel
func TimedChannelUnlock(s *discordgo.Session) {
	if config.Cfg.AutoLocker.Enabled == false {
		return
	}

	var checkInterval time.Duration = 60
	var perms int64 = 2251673408

	fmt.Println("[INIT OK] Channel unlock system module initialized")

	for {
		if time.Now().Weekday() == config.Cfg.AutoLocker.TimeDayUnlock && time.Now().Hour() == config.Cfg.AutoLocker.TimeHourUnlock && time.Now().Minute() == config.Cfg.AutoLocker.TimeMinuteUnlock {
			//Unlock the channel
			//TargetType 0 = roleID, 1 = memberID
			UnlockTrustedChannel(s, perms, TargetTypeRoleID)

			fmt.Println("[OK] Opened the channel " + config.Cfg.RoleTrusted.ChannelTrustedID)
		} else if time.Now().Weekday() == config.Cfg.AutoLocker.TimeDayLock && time.Now().Hour() == config.Cfg.AutoLocker.TimeHourLock && time.Now().Minute() == config.Cfg.AutoLocker.TimeMinuteLock {
			//Lock the channel
			//TargetType 0 = roleID, 1 = memberID
			LockTrustedChannel(s, perms, TargetTypeRoleID)
			fmt.Println("[OK] Closed the channel because regular time" + config.Cfg.RoleTrusted.ChannelTrustedID)
		} else if time.Now().Hour() == config.Cfg.AutoLocker.TimeHourLock && time.Now().Minute() == config.Cfg.AutoLocker.TimeMinuteLock && LockChannelToday == true && (time.Now().Weekday() != time.Sunday || time.Now().Weekday() != time.Saturday) {
			LockTrustedChannel(s, perms, TargetTypeRoleID)
			LockChannelToday = false
			fmt.Println("[OK] Closed the channel because special event ended " + config.Cfg.RoleTrusted.ChannelTrustedID)
		}

		time.Sleep(checkInterval * time.Second)
	}

}

// OneTimeChannelUnlock when a new trusted user is given a role, unlock the channel as a reward.
//TODO: Make this also automatically lock the channel
func OneTimeChannelUnlock(s *discordgo.Session, m *discordgo.MessageCreate) {
	var membersCached []*discordgo.Member

	for i := range ReadyInfoPublic.Guilds {
		if ReadyInfoPublic.Guilds[i].ID == m.GuildID {
			membersCached = ReadyInfoPublic.Guilds[i].Members
		}
	}
	fmt.Println(membersCached)

}

/*
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

*/

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

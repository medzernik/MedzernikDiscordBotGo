// Package responder_functions This file contains the logic for the next-generation Commands, instead of the old prefix based responses.
package responder_functions

import (
	"bufio"
	"database/sql"
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"github.com/medzernik/SlovakiaDiscordBotGo/database"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// AgeJoinedCMD Checks the age of the user on join
func AgeJoinedCMD(s *discordgo.Session, m *discordgo.InteractionCreate, cmd []interface{}) {

	//userId := m.ApplicationCommandData().Options[0].UserValue(s).ID
	userId := fmt.Sprintf("%s", cmd[0])
	fmt.Println(userId)

	//Every time a command is run, get a list of all users. This serves the purpose to then print the name of the corresponding user.
	//TODO: cache it in redis
	membersCached := GetMemberListFromGuild(s, config.Cfg.ServerInfo.GuildIDNumber)

	var userName string

	for i := range membersCached {
		if membersCached[i].User.ID == userId {
			userName = membersCached[i].User.Username
		} else if membersCached[i].User.ID != userId && membersCached[i].User.ID == "" {
			command.SendTextEmbedCommand(s, m.ChannelID, CommandStatusBot.ERR, m.Data.Type().String()+" : not a number or a mention", discordgo.EmbedTypeRich)
			return
		}
	}

	userTimeRaw, err := SnowflakeTimestamp(userId)
	if err != nil {
		command.SendTextEmbedCommand(s, m.ChannelID, CommandStatusBot.ERR, m.Data.Type().String()+" : not a number or a mention", discordgo.EmbedTypeRich)
		return
	}

	userTime := time.Now().Sub(userTimeRaw)

	roky := int64(userTime.Hours() / 24 / 365)
	dny := roky * 365 / 24
	hodiny := int64(userTime.Hours()) - dny/24
	minuty := int64(userTime.Minutes()) - int64(userTime.Hours())*60
	sekundy := int64(userTime.Seconds()) - int64(userTime.Minutes())*60

	rokyString := strconv.FormatInt(roky, 10)
	dnyString := strconv.FormatInt(dny, 10)
	hodinyString := strconv.FormatInt(hodiny, 10)
	minutyString := strconv.FormatInt(minuty, 10)
	sekundyString := strconv.FormatInt(sekundy, 10)

	//send the embed
	command.SendTextEmbedCommand(s, m.ChannelID, CommandStatusBot.OK+userName, command.ParseStringToMentionID(userId)+" "+
		" has an account age of:\n"+
		""+rokyString+" rokov\n"+
		""+dnyString+" dni\n"+
		""+hodinyString+" hodin\n"+
		""+minutyString+" minut\n"+sekundyString+" sekund"+"<:peepoLove:687313976043765810>"+
		"", discordgo.EmbedTypeRich)
	return
}

// MuteCMD Muting function
func MuteCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	//Variable initiation
	var authorisedAdmin bool = false
	var authorisedTrusted bool = false
	authorisedAdmin = command.VerifyAdminCMD(s, cmd.ChannelID, &authorisedAdmin, cmd)
	authorisedTrusted = command.VerifyTrustedCMD(s, cmd.ChannelID, &authorisedTrusted, cmd)

	timeToCheckUsers := 24.0 * -1.0

	fmt.Println("authorisedadmin", authorisedAdmin)
	fmt.Println("trusteduser", authorisedTrusted)

	//Verify, if user has any rights at all
	if authorisedAdmin == false && authorisedTrusted == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "Error muting a user - insufficient rights.", discordgo.EmbedTypeRich)
		return
	}

	//Added only after the first check of rights, to prevent spamming of the requests
	membersCached := GetMemberListFromGuild(s, config.Cfg.ServerInfo.GuildIDNumber)
	var MuteUserString []string

	MuteUserString = append(MuteUserString, command.ParseMentionToString(fmt.Sprintf("%s", m[0])))

	//Verify for the admin role before muting.
	if authorisedAdmin == true {
		for i := range membersCached {
			for j := range MuteUserString {
				if membersCached[i].User.ID == MuteUserString[j] {
					//Try to mute
					s.GuildMemberMute(config.Cfg.ServerInfo.GuildIDNumber, MuteUserString[j], true)
					err2 := s.GuildMemberRoleAdd(config.Cfg.ServerInfo.GuildIDNumber, MuteUserString[j], config.Cfg.MuteFunction.MuteRoleID)
					if err2 != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error muting a user - cannot assign the MuteRole."+
							" "+config.Cfg.MuteFunction.MuteRoleID, discordgo.EmbedTypeRich)
						return
					}
					command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"MUTED", "Muted user "+command.ParseStringToMentionID(membersCached[i].User.ID)+" (ID: "+
						""+membersCached[i].User.ID+")", discordgo.EmbedTypeRich)
					s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[LOG]** Administrator user "+cmd.Member.User.Username+" Muted user: "+
						""+command.ParseStringToMentionID(membersCached[i].User.ID))
					return
				}

			}
		}
	}

	//If not, verify for the role of Trusted to try to mute
	if authorisedTrusted == true && authorisedAdmin == false && config.Cfg.MuteFunction.TrustedMutingEnabled == true {
		for i := range membersCached {
			for j := range MuteUserString {
				userTimeJoin, _ := membersCached[i].JoinedAt.Parse()
				timevar := userTimeJoin.Sub(time.Now()).Hours()
				if membersCached[i].User.ID == MuteUserString[j] && timevar > timeToCheckUsers {
					//Error checking
					s.GuildMemberMute(config.Cfg.ServerInfo.GuildIDNumber, MuteUserString[j], true)

					err2 := s.GuildMemberRoleAdd(config.Cfg.ServerInfo.GuildIDNumber, MuteUserString[j], config.Cfg.MuteFunction.MuteRoleID)
					if err2 != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error muting a user - cannot assign the MuteRole."+
							" "+config.Cfg.MuteFunction.MuteRoleID, discordgo.EmbedTypeRich)
						return
					}
					command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"MUTED", "Muted user younger than "+
						""+strconv.FormatInt(int64(timeToCheckUsers*-1.0), 10)+MuteUserString[j], discordgo.EmbedTypeRich)

					s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[LOG]** Trusted user "+command.ParseStringToMentionID(cmd.User.Username)+" Muted user: "+
						""+command.ParseStringToMentionID(membersCached[i].User.ID))
					return

					//muting cannot be done if the time limit has been passed
				} else if membersCached[i].User.ID == MuteUserString[j] && timevar < timeToCheckUsers {
					command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "Trusted users cannot mute anyone who has joined more than "+
						""+strconv.FormatInt(int64(timeToCheckUsers*-1.0), 10)+" hours ago.", discordgo.EmbedTypeRich)
					return
				}
			}
		}

	} else if config.Cfg.MuteFunction.TrustedMutingEnabled == false && authorisedTrusted == true && authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.WARN, "Muting by Trusted users is currently disabled"+
			" "+config.Cfg.MuteFunction.MuteRoleID, discordgo.EmbedTypeRich)
		return
	} else {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "Undefined permissions error"+
			" "+config.Cfg.MuteFunction.MuteRoleID, discordgo.EmbedTypeRich)
		return
	}
	return
}

// UnmuteCMD Unmuting function
func UnmuteCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	//Variable initiation
	var authorisedAdmin bool = false
	var authorisedTrusted bool = false
	authorisedAdmin = command.VerifyAdminCMD(s, cmd.ChannelID, &authorisedAdmin, cmd)
	authorisedTrusted = command.VerifyTrustedCMD(s, cmd.ChannelID, &authorisedTrusted, cmd)

	timeToCheckUsers := 24.0 * -1.0

	//Verify, if user has any rights at all
	if authorisedAdmin == false && authorisedTrusted == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "Error unmuting a user - insufficient rights.", discordgo.EmbedTypeRich)
		return
	}

	//Added only after the first check of rights, to prevent spamming of the requests
	membersCached := GetMemberListFromGuild(s, config.Cfg.ServerInfo.GuildIDNumber)
	var UnmuteUserString []string

	UnmuteUserString = append(UnmuteUserString, command.ParseMentionToString(fmt.Sprintf("%s", m[0])))

	//Verify for the admin role before muting.
	if authorisedAdmin == true {
		for i := range membersCached {
			for j := range UnmuteUserString {
				if membersCached[i].User.ID == UnmuteUserString[j] {
					//Try to mute
					s.GuildMemberMute(config.Cfg.ServerInfo.GuildIDNumber, UnmuteUserString[j], false)
					err2 := s.GuildMemberRoleRemove(config.Cfg.ServerInfo.GuildIDNumber, UnmuteUserString[j], config.Cfg.MuteFunction.MuteRoleID)
					if err2 != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error Unmuting a user - cannot remove the MuteRole."+
							" "+config.Cfg.MuteFunction.MuteRoleID, discordgo.EmbedTypeRich)
						return
					}
					command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"UNMUTED", "Unmuted user "+command.ParseStringToMentionID(membersCached[i].User.ID)+" (ID: "+
						""+membersCached[i].User.ID+")", discordgo.EmbedTypeRich)
					s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[LOG]** Administrator user "+cmd.Member.User.Username+" Unmuted user: "+
						""+command.ParseStringToMentionID(membersCached[i].User.ID))
					return
				}

			}
		}
	}

	//If not, verify for the role of Trusted to try to mute
	if authorisedTrusted == true && authorisedAdmin == false && config.Cfg.MuteFunction.TrustedMutingEnabled == true {
		for i := range membersCached {
			for j := range UnmuteUserString {
				userTimeJoin, _ := membersCached[i].JoinedAt.Parse()
				timevar := userTimeJoin.Sub(time.Now()).Hours()
				if membersCached[i].User.ID == UnmuteUserString[j] && timevar > timeToCheckUsers {
					//Error checking
					s.GuildMemberMute(config.Cfg.ServerInfo.GuildIDNumber, UnmuteUserString[j], false)

					err2 := s.GuildMemberRoleRemove(config.Cfg.ServerInfo.GuildIDNumber, UnmuteUserString[j], config.Cfg.MuteFunction.MuteRoleID)
					if err2 != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error Unmuting a user - cannot remove the MuteRole."+
							" "+config.Cfg.MuteFunction.MuteRoleID, discordgo.EmbedTypeRich)
						return
					}
					command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"UNMUTED", "Unmuted user younger than "+
						""+strconv.FormatInt(int64(timeToCheckUsers*-1.0), 10)+UnmuteUserString[j], discordgo.EmbedTypeRich)

					s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[LOG]** Trusted user "+command.ParseStringToMentionID(cmd.User.Username)+" Unmuted user: "+
						""+command.ParseStringToMentionID(membersCached[i].User.ID))
					return

					//muting cannot be done if the time limit has been passed
				} else if membersCached[i].User.ID == UnmuteUserString[j] && timevar < timeToCheckUsers {
					command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "Trusted users cannot unmuted anyone who has joined more than "+
						""+strconv.FormatInt(int64(timeToCheckUsers*-1.0), 10)+" hours ago.", discordgo.EmbedTypeRich)
					return
				}
			}
		}

	} else if config.Cfg.MuteFunction.TrustedMutingEnabled == false && authorisedTrusted == true && authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.WARN, "Unmuting by Trusted users is currently disabled"+
			" "+config.Cfg.MuteFunction.MuteRoleID, discordgo.EmbedTypeRich)
		return
	} else {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "Undefined permissions error"+
			" "+config.Cfg.MuteFunction.MuteRoleID, discordgo.EmbedTypeRich)
		return
	}
	return
}

func KickUserCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	var reasonExists bool
	var reason string
	var authorisedAdmin bool = false
	authorisedAdmin = command.VerifyAdminCMD(s, cmd.ChannelID, &authorisedAdmin, cmd)

	if authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "Error kicking a user - insufficient rights.", discordgo.EmbedTypeRich)
		return
	}

	membersCached := GetMemberListFromGuild(s, config.Cfg.ServerInfo.GuildIDNumber)
	if len(m) > 1 {
		reason = fmt.Sprintf("%s", m[1])
		reasonExists = true
	}

	var KickUserString string = command.ParseMentionToString(fmt.Sprintf("%s", m[0]))

	s.ChannelMessageSend(cmd.ChannelID, "**[PERM]** Permissions check complete.")

	if authorisedAdmin == true {
		for i := range membersCached {
			if membersCached[i].User.ID == KickUserString {
				if reasonExists == true {
					//DM the user of his kick + reason
					userNotifChanID, err0 := s.UserChannelCreate(KickUserString)
					if err0 != nil {
						s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error notifying the user of his kick")
					} else {
						s.ChannelMessageSend(userNotifChanID.ID, "You have been kicked from the server. Reason: "+reason)
					}

					//perform the kick itself
					err := s.GuildMemberDeleteWithReason(config.Cfg.ServerInfo.GuildIDNumber, KickUserString, reason)
					if err != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error kicking user ID "+membersCached[i].User.ID, discordgo.EmbedTypeRich)
						return
					}

					//log the kick
					command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"KICKED", "Kicked user "+membersCached[i].User.Username+"for "+
						""+reason, discordgo.EmbedTypeRich)
					s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "User "+KickUserString+" ,was kicked for: "+fmt.Sprintf("%s", m[1])+" .Kicked by "+cmd.Member.Nick)
				} else {
					//DM the user of his kick
					userNotifChanID, err0 := s.UserChannelCreate(KickUserString)
					if err0 != nil {
						s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error notifying the user of his kick")
					} else {
						s.ChannelMessageSend(userNotifChanID.ID, "You have been kicked from the server.")
					}

					//perform the kick itself
					err := s.GuildMemberDelete(config.Cfg.ServerInfo.GuildIDNumber, KickUserString)
					if err != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error kicking user ID "+membersCached[i].User.ID, discordgo.EmbedTypeRich)
						return
					}

					//log the kick
					//TODO: Fix the KickUserString -> stringf(m[0])
					command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"KICKED", "Kicked user "+membersCached[i].User.Username, discordgo.EmbedTypeRich)
					s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "User "+KickUserString+" "+fmt.Sprintf("%s", m[0])+" Kicked by "+cmd.Member.Nick)
				}
			}

		}
	}
	return
}

// BanUser bans a user
func BanUserCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	var reason string
	var reasonExists bool = false
	var daysDelete int = 7
	var authorisedAdmin bool = false
	authorisedAdmin = command.VerifyAdminCMD(s, cmd.ChannelID, &authorisedAdmin, cmd)

	if authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "Error banning a user - insufficient rights.", discordgo.EmbedTypeRich)
		return
	}

	membersCached := GetMemberListFromGuild(s, config.Cfg.ServerInfo.GuildIDNumber)

	var BanUserString string = fmt.Sprintf("%s", m[0])
	if len(m) > 1 {
		reason = fmt.Sprintf("%s", m[1])
		reasonExists = true
	}

	s.ChannelMessageSend(cmd.ChannelID, "**[PERM]** Permissions check complete.")

	if authorisedAdmin == true {
		for i := range membersCached {
			if membersCached[i].User.ID == BanUserString {
				if reasonExists == true {
					userNotifChanID, err0 := s.UserChannelCreate(BanUserString)
					if err0 != nil {
						s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error notifying the user of his ban")
					} else {
						s.ChannelMessageSend(userNotifChanID.ID, "You have been banned from the server. Reason: "+reason)
					}

					err := s.GuildBanCreateWithReason(config.Cfg.ServerInfo.GuildIDNumber, BanUserString, reason, daysDelete)
					if err != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error banning user ID "+membersCached[i].User.ID, discordgo.EmbedTypeRich)
						return
					}
					command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"BANNED", "Banned user "+membersCached[i].User.Username+"for "+
						""+reason, discordgo.EmbedTypeRich)
				} else {
					userNotifChanID, err0 := s.UserChannelCreate(BanUserString)
					if err0 != nil {
						s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error notifying the user of his ban")
					} else {
						s.ChannelMessageSend(userNotifChanID.ID, "You have been banned from the server.")
					}

					err1 := s.GuildBanCreate(config.Cfg.ServerInfo.GuildIDNumber, BanUserString, daysDelete)
					if err1 != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error banning user ID "+membersCached[i].User.ID, discordgo.EmbedTypeRich)
						return
					}
					command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"BANNED", "Banning user "+membersCached[i].User.Username, discordgo.EmbedTypeRich)
					//TODO: Fix the BanUserString -> stringf(m[0])
					s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "User "+BanUserString+" "+fmt.Sprintf("%s", m[0])+" Banned by "+cmd.Member.Nick)
				}
			}

		}
	}
	return
}

// CheckUsersCMD Checks the age of users
func CheckUsersCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	var timeToCheckUsers int64
	if len(m) > 0 {
		timeToCheckUsers = m[0].(int64)
		timeToCheckUsers *= -1
	} else {
		timeToCheckUsers = 24 * -1
	}

	//variable definitions
	var authorisedAdmin bool = false
	authorisedAdmin = command.VerifyAdminCMD(s, cmd.ChannelID, &authorisedAdmin, cmd)

	if authorisedAdmin == true {
		membersCached := GetMemberListFromGuild(s, config.Cfg.ServerInfo.GuildIDNumber)
		var mainOutputMsg string
		var IDOutputMsg string

		//iterate over the members_cached array. Maximum limit is 1000.
		for i := range membersCached {
			userTimeJoin, _ := membersCached[i].JoinedAt.Parse()
			var timeVar int64 = int64(userTimeJoin.Sub(time.Now()).Hours())

			if timeVar > timeToCheckUsers {
				mainOutputMsg += "This user is too young (less than " +
					"" + strconv.FormatInt(timeToCheckUsers*-1, 10) + "h join age): " +
					"" + membersCached[i].User.Username + " ,**ID:** " +
					"" + membersCached[i].User.ID + "\n"
				IDOutputMsg += membersCached[i].User.ID + " "
			}
		}
		//print out the amount of members_cached (max is currently 1000)
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"RECENT USERS", mainOutputMsg+"\n**IDs of the users (copyfriendly):**\n"+IDOutputMsg, discordgo.EmbedTypeRich)
	} else if authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "You do not have the permission to use this command", discordgo.EmbedTypeRich)
		return
	}
	return

}

// PlanGameCMD Plans a game for a person with a timed reminder
func PlanGameCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	go GamePlanInsertCMD(s, cmd, m)
	return
}

// PlannedGamesCMD Checks the planned games and outputs them into the guild
func PlannedGamesCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	//open database and then close it (defer)
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer func(sqliteDatabase *sql.DB) {
		err := sqliteDatabase.Close()
		if err != nil {
			fmt.Println("error closing the database: ", err)
		}
	}(sqliteDatabase)

	var plannedGames string
	database.DisplayAllGamesPlanned(sqliteDatabase, &plannedGames)

	//send info to channel
	command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"PLANNED GAMES", plannedGames, discordgo.EmbedTypeRich)
	return
}

// GamePlanInsertCMD Inserts the game into the database
func GamePlanInsertCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	//open database and then close it (defer)
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer func(sqliteDatabase *sql.DB) {
		err := sqliteDatabase.Close()
		if err != nil {

		}
	}(sqliteDatabase)

	//transform to timestamp
	splitTimeArgument := strings.Split(fmt.Sprintf("%s", m[0]), ":")

	//Put hours into timeHours
	timeHour, err := strconv.Atoi(splitTimeArgument[0])
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error converting hours", discordgo.EmbedTypeRich)
		return
	}
	//put minutes into timeMinute
	timeMinute, err := strconv.Atoi(splitTimeArgument[1])
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error converting minutes", discordgo.EmbedTypeRich)
		//(*s).ChannelMessageSend((*m).ChannelID, "**[ERR]** Error converting minutes")
		//fmt.Printf("%s", err)
		return
	}
	//get current date and replace hours and minutes with user variables
	gameTimestamp := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), timeHour, timeMinute, time.Now().Second(), 0, time.Now().Location())
	gameTimestampInt := gameTimestamp.Unix()

	//export to database
	database.InsertGame(sqliteDatabase, gameTimestampInt, fmt.Sprintf("%s", m[1]), fmt.Sprintf("%s", m[2]))

	var plannedGames string
	database.DisplayGamePlanned(sqliteDatabase, &plannedGames)

	command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"PLANNED A GAME", plannedGames, discordgo.EmbedTypeRich)
	return
}

// TopicCMD Outputs a random topic for discussion found in topic_questions.txt
func TopicCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	fileHandle, err := os.Open("topic_questions.txt")
	if err != nil {
		fmt.Println("error reading the file: ", err)
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error reading the file topic_questions.txt", discordgo.EmbedTypeRich)
		return
	}
	defer func(fileHandle *os.File) {
		err := fileHandle.Close()
		if err != nil {
			fmt.Println("**[ERR]** error closing the file with topics")
		}
	}(fileHandle)

	fileScanner := bufio.NewScanner(fileHandle)

	var splitTopic []string

	for fileScanner.Scan() {
		splitTopic = append(splitTopic, fileScanner.Text())
	}

	//a, b is the length of the topic.
	a := 0
	b := len(splitTopic)

	rand.Seed(time.Now().UnixNano())
	n := a + rand.Intn(b-a+1)

	//this checks the slice length and prevents a panic if for any chance it happened. Just in case.
	if n > len(splitTopic)-1 {
		fmt.Println("**[ERR_PARSE]** Slice is smaller than allowed\n This error should not have ever happened...")
		return
	}

	command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"TOPIC", splitTopic[n], discordgo.EmbedTypeRich)
	return
}

// GetWeatherCMD outputs weather information from openWeatherMap
func GetWeatherCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	type wData struct {
		name       string
		weather    string
		condition  string
		temp       string
		tempMax    string
		tempMin    string
		tempFeel   string
		pressure   string
		humidity   string
		windSpeed  string
		rainAmount string
		sunrise    string
		sunset     string
	}

	w, err := owm.NewCurrent("C", "en", config.Cfg.ServerInfo.WeatherAPIKey)
	if err != nil {
		fmt.Println("Error processing the request")
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error processing the request", discordgo.EmbedTypeRich)
	}

	var commandString string = fmt.Sprintf("%s", m[0])

	err2 := w.CurrentByName(commandString)
	if err2 != nil {
		log.Println(err2)
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "The city "+commandString+" does not exist", discordgo.EmbedTypeRich)
		return
	}

	var weatherData = wData{
		name:       w.Name,
		weather:    w.Weather[0].Main,
		condition:  w.Weather[0].Description,
		temp:       strconv.FormatFloat(w.Main.Temp, 'f', 1, 64) + " °C",
		tempMax:    strconv.FormatFloat(w.Main.TempMax, 'f', 1, 64) + " °C",
		tempMin:    strconv.FormatFloat(w.Main.TempMin, 'f', 1, 64) + " °C",
		tempFeel:   strconv.FormatFloat(w.Main.FeelsLike, 'f', 1, 64) + " °C",
		pressure:   strconv.FormatFloat(w.Main.Pressure, 'f', 1, 64) + " hPa",
		humidity:   strconv.FormatInt(int64(w.Main.Humidity), 10) + " %",
		windSpeed:  strconv.FormatFloat(w.Wind.Speed, 'f', 1, 64) + " km/h",
		rainAmount: strconv.FormatFloat(w.Rain.OneH*10, 'f', 1, 64) + " %",
		sunrise:    time.Unix(int64(w.Sys.Sunrise), 0).Format(time.Kitchen),
		sunset:     time.Unix(int64(w.Sys.Sunset), 0).Format(time.Kitchen),
	}

	/*
		var weatherDataString string = "```\n" +
			"City:\t\t" + weatherData.name + "\n" +
			"Weather:\t" + weatherData.weather + "\n" +
			"Condition:\t" + weatherData.condition + "\n" +
			"Temperature:" + weatherData.temp + "\n" +
			"Max Temp:\t" + weatherData.tempMax + "\n" +
			"Min Temp:\t" + weatherData.tempMin + "\n" +
			"Feel Temp:\t" + weatherData.tempFeel + "\n" +
			"Pressure:\t" + weatherData.pressure + "\n" +
			"Humidity:\t" + weatherData.humidity + "\n" +
			"Wind Speed:\t" + weatherData.windSpeed + "\n" +
			"Rainfall:\t" + weatherData.rainAmount + "\n" +
			"Sunrise:\t" + weatherData.sunrise + "\n" +
			"Sunset:\t" + weatherData.sunset + "\n" +
			"```"

	*/

	embed := NewEmbed().
		SetTitle("WEATHER IN: "+strings.ToUpper(weatherData.name)).
		SetDescription(":cloud: **"+strings.ToUpper(weatherData.condition)+"**").
		AddField("**TEMPERATURE**",
			":thermometer: "+weatherData.temp+
				"\n:high_brightness: "+weatherData.tempMax+
				"\n:low_brightness: "+weatherData.tempMin).
		AddField("**PRESSURE**", ":diving_mask: "+weatherData.pressure).
		AddField("**HUMIDITY**", ":sweat: "+weatherData.humidity).
		AddField("**WIND SPEED**", ":dash: "+weatherData.windSpeed).
		AddField("**RAINFALL**", ":droplet: "+weatherData.rainAmount).
		AddField("**SUN RISE/SET**",
			":sunrise_over_mountains: "+weatherData.sunrise+
				"\n:city_sunset: "+weatherData.sunset).
		InlineAllFields().
		SetColor(3066993).MessageEmbed

	s.ChannelMessageSendEmbed(cmd.ChannelID, embed)
	return

}

// PurgeMessagesCMD Purges messages 1-100 in the current channel
func PurgeMessagesCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	var authorisedAdmin bool = false
	authorisedAdmin = command.VerifyAdminCMD(s, cmd.ChannelID, &authorisedAdmin, cmd)

	if authorisedAdmin == true {
		var messageArrayToDelete []string

		/*
			numMessages, err1 := strconv.ParseInt(fmt.Sprintf("%s", m[0]), 10, 64)
			if err1 != nil {
				command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Invalid number provided", discordgo.EmbedTypeRich)
				return
			}

		*/

		numMessages := m[0].(int64)
		if numMessages > 99 || numMessages < 1 {
			command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.SYNTAX, "The min-max of the number is 1-100", discordgo.EmbedTypeRich)
			return
		}

		messageArrayComplete, err1 := s.ChannelMessages(cmd.ChannelID, int(numMessages), cmd.ID, "", "")
		if err1 != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Cannot get the ID of messages", discordgo.EmbedTypeRich)
			return
		}

		for i := range messageArrayComplete {
			messageArrayToDelete = append(messageArrayToDelete, messageArrayComplete[i].ID)
		}

		err2 := s.ChannelMessagesBulkDelete(cmd.ChannelID, messageArrayToDelete)
		if err2 != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error deleting the requested messages...", discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"PURGED", "Purged "+strconv.FormatInt(int64(len(messageArrayToDelete)), 10)+" "+
			"messages", discordgo.EmbedTypeRich)

		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[LOG]** User "+cmd.Member.Nick+" deleted "+strconv.FormatInt(int64(len(messageArrayToDelete)), 10)+" messages in channel "+"<#"+cmd.ChannelID+">")

		return
	} else {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "Insufficient permissions.", discordgo.EmbedTypeRich)
		return
	}

}

// MembersCMD outputs the number of current members of the server. No returns
func MembersCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	memberList := GetMemberListFromGuild(s, config.Cfg.ServerInfo.GuildIDNumber)
	memberListLength := uint64(len(memberList))

	command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+strconv.FormatUint(memberListLength, 10), ""+
		"There are "+strconv.FormatUint(memberListLength, 10)+" members on the server", discordgo.EmbedTypeRich)

	return
}

// PruneCountCMD outputs the number of users that could be pruned
func PruneCountCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	pruneDaysInt := m[0].(int64)

	if pruneDaysInt < 7 || pruneDaysInt > 30 {
		pruneDaysInt = 0
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.WARN, "Command is limited to range 7-30 for safety reasons", discordgo.EmbedTypeRich)
		return
	}

	pruneDaysCount, err := s.GuildPruneCount(config.Cfg.ServerInfo.GuildIDNumber, uint32(pruneDaysInt))
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error checking members to prune.", discordgo.EmbedTypeRich)
		return
	}

	command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+strconv.FormatUint(uint64(pruneDaysCount), 10), ""+
		"There are "+strconv.FormatUint(uint64(pruneDaysCount), 10)+" members to prune", discordgo.EmbedTypeRich)

	return
}

// PruneMembersCMD prunes members
func PruneMembersCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	var authorisedAdmin bool = false
	authorisedAdmin = command.VerifyAdminCMD(s, cmd.ChannelID, &authorisedAdmin, cmd)

	if authorisedAdmin == true {
		//request prune number amount
		pruneDaysCountInt := m[0].(int64)

		var pruneDaysCountUInt = uint32(pruneDaysCountInt)

		if pruneDaysCountInt == 0 {
			command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.SYNTAX, "Cannot prune time of 0 days. Allowed frame is 7-30", discordgo.EmbedTypeRich)
			s.ChannelMessageSend(cmd.ChannelID, "**[ERR]** Invalid days to prune (0)")
			return
		}

		//prunes the members and assigns the result of the pruned members count to a variable
		prunedMembersCount, err1 := s.GuildPrune(config.Cfg.ServerInfo.GuildIDNumber, pruneDaysCountUInt)
		if err1 != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error pruning members", discordgo.EmbedTypeRich)
		}

		//log output

		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"PRUNED", strconv.FormatInt(int64(prunedMembersCount), 10)+
			" members from the server", discordgo.EmbedTypeRich)
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[LOG]** User "+cmd.Member.Nick+
			" used a prune and kicked "+strconv.FormatInt(int64(prunedMembersCount), 10)+" members")
		return

		//permission output
	} else {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "Insufficient permissions", discordgo.EmbedTypeRich)
		return
	}

}

// SetRoleChannelPerm sets a channel permission using an int value
func SetRoleChannelPermCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	var verifyAdmin bool = false
	verifyAdmin = command.VerifyAdminCMD(s, cmd.ChannelID, &verifyAdmin, cmd)

	//check if user is admin before using the command
	if verifyAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "You do not have the permission to change the channel permissions for a role.", discordgo.EmbedTypeRich)
		return
	}

	//get the role set
	permissionRole := command.ParseRoleMentionToString(fmt.Sprintf("%s", m[1]))
	//get the ID set
	permissionID := m[2].(int64)
	//get whether to allow or deny the permissions
	permissionAllow := m[0].(bool)

	// GET WHETHER TO SET PERMS FOR A ROLE OR A MEMBER

	if permissionAllow == true {
		err := s.ChannelPermissionSet(cmd.ChannelID, permissionRole, 0, permissionID, 0)
		if err != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error setting the permissions on the channel."+err.Error(), discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"PERMISSIONS ALLOWED", "Permissions "+strconv.FormatInt(permissionID, 10)+" successfully allowed"+
			"", discordgo.EmbedTypeRich)
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Admin "+cmd.Member.Nick+" denied permissions "+strconv.FormatInt(permissionID, 10)+""+
			" for channel "+command.ParseStringToChannelID(cmd.ChannelID)+" to a role "+command.ParseStringToRoleMention(permissionRole))
		return

		//if not allowed then not allow
	} else if permissionAllow == false {
		err := s.ChannelPermissionSet(cmd.ChannelID, permissionRole, 0, 0, permissionID)
		if err != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error setting the permissions on the channel."+err.Error(), discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"PERMISSIONS DENIED", "Permissions "+strconv.FormatInt(permissionID, 10)+" successfully denied"+
			"", discordgo.EmbedTypeRich)
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Admin "+cmd.Member.Nick+" denied permissions "+strconv.FormatInt(permissionID, 10)+""+
			" for channel "+command.ParseStringToChannelID(cmd.ChannelID)+" to a role "+command.ParseStringToRoleMention(permissionRole))
		return

		//if there is an invalid syntax with the allow/deny argument then
	} else {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.SYNTAX, "Invalid syntax, .setpermission **<allow,deny>** "+
			"@roletoset INTPERMIDS", discordgo.EmbedTypeRich)
		return
	}

}

// SetUserChannelPermCMD sets a channel permission using an int value
func SetUserChannelPermCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	var verifyAdmin bool = false
	verifyAdmin = command.VerifyAdminCMD(s, cmd.ChannelID, &verifyAdmin, cmd)

	//check if user is admin before using the command
	if verifyAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "You do not have the permission to change the channel permissions for a role.", discordgo.EmbedTypeRich)
		return
	}

	//get the role to set
	permissionRole := command.ParseRoleMentionToString(fmt.Sprintf("%s", m[1]))

	permissionID := m[2].(int64)

	//get whether to allow or deny the permissions
	permissionAllow := m[0].(bool)
	if permissionAllow == true {
		err := s.ChannelPermissionSet(cmd.ChannelID, permissionRole, 1, permissionID, 0)
		if err != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error setting the permissions on the channel."+err.Error(), discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"PERMISSIONS ALLOWED", "Permissions "+strconv.FormatInt(permissionID, 10)+" successfully allowed"+
			"", discordgo.EmbedTypeRich)
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Admin "+cmd.ChannelID+" denied permissions "+strconv.FormatInt(permissionID, 10)+""+
			" for channel "+command.ParseStringToChannelID(cmd.ChannelID)+" to a user "+command.ParseStringToMentionID(permissionRole))
		return

		//if not allowed then not allow
	} else if permissionAllow == false {
		err := s.ChannelPermissionSet(cmd.ChannelID, permissionRole, 1, 0, permissionID)
		if err != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error setting the permissions on the channel."+err.Error(), discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"PERMISSIONS DENIED", "Permissions "+strconv.FormatInt(permissionID, 10)+" successfully denied"+
			"", discordgo.EmbedTypeRich)
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Admin "+cmd.Member.Nick+" denied permissions "+strconv.FormatInt(permissionID, 10)+""+
			" for channel "+command.ParseStringToChannelID(cmd.ChannelID)+" to a user "+command.ParseStringToMentionID(permissionRole))
		return
	}
}

//RedirectDiscussionCMD  sets a channel to a big slowmode for 10 minutes and then redirects the conversation elsewhere. When threads become available, sets the thread and more...
func RedirectDiscussionCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	var verifyAdmin bool = false
	verifyAdmin = command.VerifyAdminCMD(s, cmd.ChannelID, &verifyAdmin, cmd)

	//check if user is admin before using the command
	if verifyAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "You do not have the permission to change the channel slowmode.", discordgo.EmbedTypeRich)
		return
	}

	originalChannelInfo, err := s.Channel(cmd.ChannelID)
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Cannot get info of the channel to modify, aborting.", discordgo.EmbedTypeRich)
		return
	}

	var channelIDString string
	var slowmodeChannelSet discordgo.ChannelEdit = discordgo.ChannelEdit{
		Name:                 originalChannelInfo.Name,
		Topic:                originalChannelInfo.Topic,
		NSFW:                 originalChannelInfo.NSFW,
		Position:             originalChannelInfo.Position,
		Bitrate:              originalChannelInfo.Bitrate,
		UserLimit:            originalChannelInfo.UserLimit,
		PermissionOverwrites: originalChannelInfo.PermissionOverwrites,
		ParentID:             originalChannelInfo.ParentID,
		RateLimitPerUser:     360,
	}

	channelIDString = command.ParseChannelToString(fmt.Sprintf("%v", m[0]))

	_, err1 := s.ChannelEditComplex(cmd.ChannelID, &slowmodeChannelSet)
	if err1 != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error setting the slowmode.", discordgo.EmbedTypeRich)
		return
	}
	command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"SLOWMODE SET FOR "+strconv.FormatInt(int64(slowmodeChannelSet.RateLimitPerUser), 10)+" SECONDS", "Continue discussion in "+
		""+command.ParseStringToChannelID(channelIDString), discordgo.EmbedTypeRich)

	return
}

// SlowModeChannelCMD sets a channel to a desired slowmode. 0 is a bugged value, so at least sets to 1 second.
func SlowModeChannelCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	var verifyAdmin bool = false
	verifyAdmin = command.VerifyAdminCMD(s, cmd.ChannelID, &verifyAdmin, cmd)

	//check if user is admin before using the command
	if verifyAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTH, "You do not have the permission to change the channel slowmode.", discordgo.EmbedTypeRich)
		return
	}

	//parse the argument to input
	numOfSeconds := m[0].(uint64)

	//verify inputs
	if numOfSeconds > 21600 {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Seconds must be in the valid range (0-21600)", discordgo.EmbedTypeRich)
		return
	} else if numOfSeconds == 0 {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.AUTOFIX, "Due to the bug in discord, setting the number to 1 (smallest possible wait time)", discordgo.EmbedTypeRich)
		numOfSeconds = 0
	}

	//get the original channel info
	originalChannelInfo, err := s.Channel(cmd.ChannelID)
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Cannot get info of the channel to modify, aborting.", discordgo.EmbedTypeRich)
		return
	}

	//get the current channel info and only modify the slowmode var
	var slowmodeChannelSet discordgo.ChannelEdit = discordgo.ChannelEdit{
		Name:                 originalChannelInfo.Name,
		Topic:                originalChannelInfo.Topic,
		NSFW:                 originalChannelInfo.NSFW,
		Position:             originalChannelInfo.Position,
		Bitrate:              originalChannelInfo.Bitrate,
		UserLimit:            originalChannelInfo.UserLimit,
		PermissionOverwrites: originalChannelInfo.PermissionOverwrites,
		ParentID:             originalChannelInfo.ParentID,
		RateLimitPerUser:     int(numOfSeconds),
	}

	//set the slowmode
	channelEdited, err1 := s.ChannelEditComplex(cmd.ChannelID, &slowmodeChannelSet)
	if err1 != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.ERR, "Error setting the slowmode.", discordgo.EmbedTypeRich)
		return
	}

	//send the confirmation message
	command.SendTextEmbedCommand(s, cmd.ChannelID, CommandStatusBot.OK+"SLOWMODE "+strconv.FormatUint(numOfSeconds, 10)+""+
		" SECONDS", "Set the channel slowmode to "+strconv.FormatUint(numOfSeconds, 10)+""+
		" seconds per message", discordgo.EmbedTypeRich)

	s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Admin "+cmd.Member.Nick+" set a "+strconv.FormatInt(int64(channelEdited.RateLimitPerUser), 10)+""+
		" seconds slowmode in channel "+command.ParseStringToChannelID(channelEdited.ID))
	return
}

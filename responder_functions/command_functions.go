// Package responder_functions This file contains the logic for the next-generation Commands, instead of the old prefix based responses.
package responder_functions

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"strconv"
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

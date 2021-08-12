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

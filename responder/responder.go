// Package responder defines the trigger words for the chat from users and then runs the appropriate commands.
package responder

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"github.com/medzernik/SlovakiaDiscordBotGo/database"
	"github.com/medzernik/SlovakiaDiscordBotGo/responder_functions"
	"strconv"
)

func RegisterPlugin(s *discordgo.Session) {
	s.AddHandler(messageCreated)
	s.AddHandler(ready)
	s.AddHandler(userUpdate)

}

// Function automatically unlocks the trusted channel until 6AM next day, when a new user becomes trusted
func userUpdate(s *discordgo.Session, m *discordgo.GuildMemberUpdate) {
	//if the autolocker is disabled, just don't do anything.
	//TODO: only check on the current role update...
	if config.Cfg.AutoLocker.AutoUnlockTrustedID1 == true {
		//set the permissions we want to set when autounlocking (calculated online on discord permissions calculator)
		var perms int64 = 2251673408
		var membersCached []*discordgo.Member

		for i := range responder_functions.ReadyInfoPublic.Guilds {
			if responder_functions.ReadyInfoPublic.Guilds[i].ID == m.GuildID {
				membersCached = responder_functions.ReadyInfoPublic.Guilds[i].Members
			}
		}
		target := responder_functions.TargetTypeRoleID

		fmt.Println("**[ROLE_UPDATE]** Checking the update role of a user " + m.Member.Nick)
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[NEW DEBILKO]** Unlocking the channel for roles and will autolock at 6 AM next day.")
		//search the roles of a user if he got the one we want.
		for i := range m.Roles {
			for j := range membersCached {
				if m.Roles[i] == config.Cfg.RoleTrusted.RoleTrustedID1 && m.Roles[j] != membersCached[j].Roles[i] {
					responder_functions.UnlockTrustedChannel(s, perms, target)
					responder_functions.LockChannelToday = true
					break
				}
			}
		}
	}

	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreated(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	cmd, err := command.ParseCommand(m.Content)

	if err != nil {
		println(err.Error())
		return
	}

	//TODO: finish the help system
	if command.IsCommand(&cmd, "help") {
		command.SendTextEmbed(s, m, responder_functions.CommandStatusBot.OK+"**The available commands are:**", "**.zasielkovna** - AAAAAA.\n"+
			"**.age @user** - checks the age of a user's account.\n"+
			"**.mute @user** - **[ADMIN] [TRUSTED]** gives a user a muted role. Works on users that joined less than 24 hours ago for [TRUSTED].\n"+
			"**.checkusers** - **[ADMIN]** lists users that joined less than 24h ago.\n"+
			"**.plan HH:MM Game_Name @user** - plans a game for a given time with a user. Bot reminds with a ping when time is met.\n"+
			"**.planned** - lists all the planned games.\n"+
			"**.topic** - outputs a random topic for discussion.\n"+
			"**.weather city name** - outputs weather information for a given city.\n"+
			"**.kick @user reason for the kick** - **[ADMIN]** kicks a user with a provided reason.\n"+
			"**.ban @user reason for the ban** - **[ADMIN]** bans a user with a provided reason and deletes 7 days of messages.\n"+
			"**.purge NUMBER (1-100)** - **[ADMIN]** purges the amount of messages in the current channel.\n"+
			"**.version** - displays the current bot version.\n"+
			"**.prunecount NUMBER (7-X)** - shows the number of pruneable (inactive) members. NUMBER = days (7 minimum).\n"+
			"**.prunemember NUMBER (7-X) ** - **[ADMIN]** prunes members inactive for NUMBER days (7 minimum).\n"+
			"**.members** - outputs the number of members on the server.\n"+
			"**.configreload** - **[ADMIN]** reloads the config file from disk.\n"+
			"**.slow NUMBER (0-21600)** - **[ADMIN]** sets a channel slowmode.\n"+
			"**.redirect #CHANNEL** - **[ADMIN]** redirects discussion in #CHANNELNAME. When Threads become available, redirect to a thread instead (TBD).\n", discordgo.EmbedTypeRich)
	}

	//version of the bot running
	if command.IsCommand(&cmd, "version") {
		command.SendTextEmbed(s, m, responder_functions.CommandStatusBot.OK+responder_functions.Version, "Version number is: "+
			""+responder_functions.Version+"\n"+
			"Feature name: "+responder_functions.VersionFeatureName, discordgo.EmbedTypeRich)
	}
	//counts the users on the server
	if command.IsCommand(&cmd, "members") {
		membersCountInt := responder_functions.Members(s, cmd, m)

		membersCountString := strconv.FormatInt(int64(membersCountInt), 10)

		go command.SendTextEmbed(s, m, responder_functions.CommandStatusBot.OK+membersCountString, "There are "+membersCountString+" members", discordgo.EmbedTypeRich)
	}
	//counts pruneable members on server
	if command.IsCommand(&cmd, "prunecount") {
		pruneDaysCount := responder_functions.PruneCount(s, cmd, m)

		switch pruneDaysCount {
		case 0:
			command.SendTextEmbed(s, m, responder_functions.CommandStatusBot.ERR, "No pruneable members exist", discordgo.EmbedTypeRich)
		default:
			command.SendTextEmbed(s, m, responder_functions.CommandStatusBot.OK+strconv.FormatInt(int64(pruneDaysCount), 10), "There are **"+strconv.FormatInt(int64(pruneDaysCount), 10)+"** members to prune", discordgo.EmbedTypeRich)
		}
	}

	//testing command
	if command.IsCommand(&cmd, "test") {
		go command.SendTextEmbed(s, m, ":bangbang: TEST", "<@206720832695828480>", discordgo.EmbedTypeRich)
	}
	if command.IsCommand(&cmd, "configreload") {
		responder_functions.ConfigurationReload(s, cmd, m)

	}

}

//Ready runs when the bot starts. Starts the automatic functions and sets the status of the bot
func ready(s *discordgo.Session, ready *discordgo.Ready) {
	//Set the status
	err := s.UpdateGameStatus(0, config.Cfg.ServerInfo.BotStatus)
	if err != nil {
		fmt.Println("error setting the bot status")
		return
	}
	//run the parallel functions
	go database.DatabaseOpen()
	go database.CheckPlannedGames(&s)
	go responder_functions.TimedChannelUnlock(s)
	responder_functions.Ready(s, ready)

}

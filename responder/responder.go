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

	//Zasielkovna EasterEgg
	if command.IsCommand(&cmd, "zasielkovna") {
		go responder_functions.Zasielkovna(s, cmd, m)
	}

	//Time since joined command
	if command.IsCommand(&cmd, "age") {
		go responder_functions.AgeJoined(s, cmd, m)
	}
	//Function that mutes a user (assigns him a muted role).
	if command.IsCommand(&cmd, "mute") {
		go responder_functions.Mute(s, cmd, m)
	}

	//right now this command checks for any 1000 users on the guild that have a join time less than 24hours, then prints the names one by one.
	if command.IsCommand(&cmd, "checkusers") {
		go responder_functions.CheckUsers(s, cmd, m)
	}
	//lets you play a game with a mention
	if command.IsCommand(&cmd, "plan") {
		go responder_functions.PlanGame(s, cmd, m)
	}
	//shows the currently planned games in the database
	if command.IsCommand(&cmd, "planned") {
		go responder_functions.PlannedGames(s, cmd, m)
	}
	// outputs a random topic for a discussion
	if command.IsCommand(&cmd, "topic") {
		go responder_functions.Topic(s, cmd, m)
	}
	//aaaaaa
	if command.IsCommand(&cmd, "fox") || command.IsCommand(&cmd, "shake") {
		go responder_functions.Fox(s, m)
	}
	//outputs a weather from openWeatherMap
	if command.IsCommand(&cmd, "weather") || command.IsCommand(&cmd, "pocasie") {
		go responder_functions.GetWeather(s, cmd, m)
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
	//kicks a user
	if command.IsCommand(&cmd, "kick") {
		go responder_functions.KickUser(s, cmd, m)
	}
	//bans a user
	if command.IsCommand(&cmd, "ban") {
		go responder_functions.BanUser(s, cmd, m)
	}
	//purges messages from a channel
	if command.IsCommand(&cmd, "purge") {
		go responder_functions.PurgeMessages(s, cmd, m)
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
	//prunes members from server
	if command.IsCommand(&cmd, "prunemembers") {
		go responder_functions.PruneMembers(s, cmd, m)
	}

	//testing command
	if command.IsCommand(&cmd, "test") {
		go command.SendTextEmbed(s, m, ":bangbang: TEST", "<@206720832695828480>", discordgo.EmbedTypeRich)
	}
	if command.IsCommand(&cmd, "configreload") {
		responder_functions.ConfigurationReload(s, cmd, m)

	}
	if command.IsCommand(&cmd, "slow") {
		go responder_functions.SlowModeChannel(s, cmd, m)
	}
	//Function that unmutes a user (removes the muted role).
	if command.IsCommand(&cmd, "unmute") {
		go responder_functions.Unmute(s, cmd, m)
	}
	if command.IsCommand(&cmd, "redirect") {
		go responder_functions.RedirectDiscussion(s, cmd, m)
	}
	if command.IsCommand(&cmd, "setchannelperm") {
		go responder_functions.SetRoleChannelPerm(s, cmd, m)
	}

}

//Ready runs when the bot starts. Starts the automatic functions and sets the status of the bot
func ready(s *discordgo.Session, _ *discordgo.Ready) {
	//Set the status
	err := s.UpdateGameStatus(0, config.Cfg.ServerInfo.BotStatus)
	if err != nil {
		fmt.Println("error setting the bot status")
		return
	}
	//run the parallel functions
	go responder_functions.CheckRegularSpamAttack(s)
	go database.Databaserun()
	go database.CheckPlannedGames(&s)
	go responder_functions.TimedChannelUnlock(s)
	go responder_functions.CommandLoop(s)

}

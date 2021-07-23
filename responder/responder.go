// Package responder defines the trigger words for the chat from users and then runs the appropriate commands.
package responder

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"github.com/medzernik/SlovakiaDiscordBotGo/database"
	"github.com/medzernik/SlovakiaDiscordBotGo/responder_functions"
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
	if command.IsCommand(&cmd, "Zasielkovna") {
		responder_functions.Zasielkovna(s, cmd, m)
	}

	//Time since joined command
	if command.IsCommand(&cmd, "age") {
		responder_functions.AgeJoined(s, cmd, m)
	}
	//Function that mutes a user (assigns him a muted role).
	if command.IsCommand(&cmd, "mute") {
		responder_functions.Mute(s, cmd, m)
	}

	//right now this command checks for any 1000 users on the guild that have a join time less than 24hours, then prints the names one by one.
	if command.IsCommand(&cmd, "checkusers") {
		responder_functions.CheckUsers(s, cmd, m)
	}
	//lets you play a game with a mention
	if command.IsCommand(&cmd, "plan") {
		responder_functions.PlanGame(s, cmd, m)
	}
	//shows the currently planned games in the database
	if command.IsCommand(&cmd, "planned") {
		responder_functions.PlannedGames(s, cmd, m)
	}
	// outputs a random topic for a discussion
	if command.IsCommand(&cmd, "topic") {
		responder_functions.Topic(s, cmd, m)
	}
	//aaaaaa
	if command.IsCommand(&cmd, "fox") || command.IsCommand(&cmd, "shake") {
		responder_functions.Fox(s, m)
	}
	//outputs a weather from openWeatherMap
	if command.IsCommand(&cmd, "weather") {
		go responder_functions.GetWeather(s, cmd, m)
	}
	//TODO: finish the help system
	if command.IsCommand(&cmd, "help") {
		s.ChannelMessageSend(m.ChannelID, "**The available commands are:**\n"+
			"**.Zasielkovna** - AAAAAA\n"+
			"**.age @user** - checks the age of a user's account\n"+
			"**.mute @user** - **[ADMIN] [TRUSTED]** gives a user a muted role. Works on users that joined less than 24 hours ago for [TRUSTED]\n"+
			"**.checkusers** - **[ADMIN]** lists users that joined less than 24h ago\n"+
			"**.plan HH:MM Game_Name @user** - plans a game for a given time with a user. Bot reminds with a ping when time is met.\n"+
			"**.planned** - lists all the planned games\n"+
			"**.topic** - outputs a random topic for discussion.\n"+
			"**.weather city name** - outputs weather information for a given city.\n"+
			"**.kick @user reason for the kick** - **[ADMIN]** kicks a user with a provided reason\n"+
			"**.ban @user reason for the ban** - **[ADMIN]** bans a user with a provided reason and deletes 7 days of messages\n"+
			"**.purge NUMBER (1-100)** - **[ADMIN]** purges the amount of messages in the current channel.\n"+
			"**.version** - displays the current bot version")
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
		s.ChannelMessageSend(m.ChannelID, "**[VERSION]** "+responder_functions.Version)
	}

}

func ready(s *discordgo.Session, _ *discordgo.Ready) {
	//set the status
	err := s.UpdateGameStatus(0, "Nove features mame aaaaaa")
	if err != nil {
		fmt.Println("error setting the bot status")
		return
	}
	//run the parallel functions
	go responder_functions.CheckRegularSpamAttack(s)
	go database.Databaserun()
	go database.CheckPlannedGames(&s)
	go responder_functions.TimedChannelUnlock(s)

}

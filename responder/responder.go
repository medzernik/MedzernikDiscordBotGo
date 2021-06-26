package responder

import (
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"github.com/medzernik/SlovakiaDiscordBotGo/database"
	"github.com/medzernik/SlovakiaDiscordBotGo/responder_functions"
)

func RegisterPlugin(s *discordgo.Session) {
	s.AddHandler(messageCreated)
	//s.AddHandler(reactionAdded)
	s.AddHandler(ready)

}

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

	//Zasielkovna EasterEgg
	if command.IsCommand(&cmd, "Zasielkovna") {
		responder_functions.Zasielkovna(s, cmd, m)
	}

	//Time since joined command
	if command.IsCommand(&cmd, "age") {
		responder_functions.AgeJoined(s, cmd, m)
	}

	if command.IsCommand(&cmd, "mute") {
		responder_functions.Mute(s, cmd, m, err)
	}

	//right now this command checks for any 1000 users on the guild that have a join time less than 24hours, then prints the names one by one.
	if command.IsCommand(&cmd, "check-users") {
		responder_functions.CheckUsers(s, cmd, m)
	}

	if command.IsCommand(&cmd, "plan") {
		responder_functions.PlanGame(s, cmd, m)
	}

	if command.IsCommand(&cmd, "planned") {
		responder_functions.PlannedGames(s, cmd, m)
	}

}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	//set the status
	err := s.UpdateGameStatus(0, "Gde mozog")
	if err != nil {
		return
	}
	//run the raid checker function
	go responder_functions.CheckRegularSpamAttack(s)
	go database.Databaserun()
	go database.CheckPlannedGames(&s)

}

//var db, err = sql.Open("sqlite", dsnURI)

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

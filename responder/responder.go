// Package responder is a utility package for engaging the basic functionality.
package responder

import (
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"github.com/medzernik/SlovakiaDiscordBotGo/logging"
	"github.com/medzernik/SlovakiaDiscordBotGo/responder_functions"
)

// RegisterPlugin registering the handlers
func RegisterPlugin(s *discordgo.Session) {
	s.AddHandler(ready)

}

//This is a primitive checker for running the basic command initialization the first time. It might have been broken.
var ranInit bool = false

//Ready runs when the bot starts. Starts the automatic functions and sets the status of the bot
func ready(s *discordgo.Session, ready *discordgo.Ready) {
	//Set the bot status according to the config file
	err := s.UpdateGameStatus(0, config.Cfg.ServerInfo.BotStatus)
	if err != nil {
		logging.Log.Errorln("error setting the bot status: ", err)
	}

	//run the parallel functions, only those enabled in the config file
	/*
		if config.Cfg.Modules.Planning == true {
			go database.DatabaseOpen()
			go database.CheckPlannedGames(&s)
		}

	*/

	//Initialize the commands for Discord, but only once... thinking how to redo this.
	if ranInit == false {
		ranInit = responder_functions.Ready(s, ready)
	}

}

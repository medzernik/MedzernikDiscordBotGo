//A Discord bot written in Go. Currently has various functions. Some are automatic, some react to chat messages.
package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"github.com/medzernik/SlovakiaDiscordBotGo/logging"
	"github.com/medzernik/SlovakiaDiscordBotGo/responder"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Initialize the config
	config.LoadConfig()
	logging.Log.Traceln("Loaded the config file.")

	//Initialize the logging system
	errLogging := logging.StartLogging()
	// If the log can't be created, use stdout.
	if errLogging != nil {
		logrus.Errorln("Failed to log to file, using default stderr")
	} else {
		logging.Log.Traceln("Loaded the logging system")
	}

	logging.Log.Infoln("--------------\nStarting the bot...")
	// Create a new Discord session using the provided bot token. Panic if failed.
	dg, err := discordgo.New("Bot " + config.Cfg.ServerInfo.ServerToken)
	if err != nil {
		logging.Log.Panicln("error creating Discord session: ", err)
	}

	//Get the intents that are needed
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	dg.Identify.Token = config.Cfg.ServerInfo.ServerToken
	dg.Identify.LargeThreshold = 250

	responder.RegisterPlugin(dg)

	// Open a websocket connection to Discord and begin listening. Panic if failed.
	err = dg.Open()
	if err != nil {
		logging.Log.Panicln("Error opening the websocket!: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	KillSignal := make(chan os.Signal, 1)
	signal.Notify(KillSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-KillSignal

	// Cleanly close down the Discord session.

	err2 := dg.Close()
	if err2 != nil {
		logging.Log.Panicln("Error closing the session: ", err2)
	}

}

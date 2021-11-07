//A Discord bot written in Go. Currently has various functions. Some are automatic, some react to chat messages.
package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"github.com/medzernik/SlovakiaDiscordBotGo/logging"
	"github.com/medzernik/SlovakiaDiscordBotGo/responder"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//Initialize the config
	config.LoadConfig()
	logging.StartLogging()
	//Logging stuff

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + config.Cfg.ServerInfo.ServerToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	//TODO: Intents still dont work well...
	//dg.Identify.Intents = discordgo.IntentsGuildMembers | discordgo.IntentsGuilds | discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsDirectMessages | discordgo.IntentsAll
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	dg.Identify.Token = config.Cfg.ServerInfo.ServerToken
	dg.Identify.LargeThreshold = 250

	responder.RegisterPlugin(dg)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	err2 := dg.Close()
	if err2 != nil {
		fmt.Println("error closing the session", err)
		return
	}

}

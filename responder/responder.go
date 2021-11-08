// Package responder is a utility package for engaging the basic functionality.
package responder

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"github.com/medzernik/SlovakiaDiscordBotGo/database"
	"github.com/medzernik/SlovakiaDiscordBotGo/logging"
	"github.com/medzernik/SlovakiaDiscordBotGo/responder_functions"
)

// RegisterPlugin registering the handlers
func RegisterPlugin(s *discordgo.Session) {
	s.AddHandler(ready)
	s.AddHandler(userUpdate)
}

// Function automatically unlocks the trusted channel until 6AM next day, when a new user becomes trusted. Deprecated currently.
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
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[NEW TRUSTED]** Unlocking the channel for roles and will autolock at 6 AM next day.")
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
	if config.Cfg.Modules.Planning == true {
		go database.DatabaseOpen()
		go database.CheckPlannedGames(&s)
	}

	if config.Cfg.Modules.TimedChannelUnlock == true {
		go responder_functions.TimedChannelUnlock(s)
	}

	//Initialize the commands for Discord, but only once... thinking how to redo this.
	if ranInit == false {
		ranInit = responder_functions.Ready(s, ready)
	}

}

// Package responder defines the trigger words for the chat from users and then runs the appropriate commands.
package responder

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"github.com/medzernik/SlovakiaDiscordBotGo/database"
	"github.com/medzernik/SlovakiaDiscordBotGo/responder_functions"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func RegisterPlugin(s *discordgo.Session) {
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

func PasswordLottery(s *discordgo.Session) {
	var stringResultCache string

	for {
		//Check the time sensitive info or skip it
		if (time.Now().Weekday() == config.Cfg.LotteryChecker.TimeDayStart && time.Now().Hour() >= config.Cfg.LotteryChecker.TimeHourStart && time.Now().Hour() <= config.Cfg.LotteryChecker.TimeHourEnd && time.Now().Minute() >= config.Cfg.LotteryChecker.TimeMinuteStart && time.Now().Minute() <= config.Cfg.LotteryChecker.TimeMinuteEnd) || config.Cfg.LotteryChecker.Enabled == false {

			response, err := http.Get("https://www.ockovacie.info/")
			if err != nil {
				fmt.Println(err)
			}
			defer response.Body.Close()
			if response.StatusCode == 200 {
				bodyText, err := ioutil.ReadAll(response.Body)
				if err != nil {
					fmt.Println(err)
				}
				bodyTextProcessed := fmt.Sprintf("%s\n", bodyText)

				stringResult := strings.SplitAfter(bodyTextProcessed, "text-align:center;white-space:pre-wrap;\">")
				stringResult = strings.Split(stringResult[2], "</h2></div></div></div></div></div>")
				fmt.Println(stringResult[0])

				if stringResult[0] != stringResultCache {
					command.SendTextEmbedCommand(s, config.Cfg.ChannelLog.GamePlannedLog, responder_functions.CommandStatusBot.OK+"HESLO DO LOTÉRIE", stringResult[0], discordgo.EmbedTypeRich)
					stringResultCache = stringResult[0]
				}

			}
		}
		time.Sleep(180 * time.Second)

	}

	//https://www.ockovacie.info/
}

//Ready runs when the bot starts. Starts the automatic functions and sets the status of the bot
func ready(s *discordgo.Session, ready *discordgo.Ready) {
	//Set the bot status according to the config file
	err := s.UpdateGameStatus(0, config.Cfg.ServerInfo.BotStatus)
	if err != nil {
		fmt.Println("error setting the bot status")
		return
	}

	//run the parallel functions only if enabled in the config file
	if config.Cfg.Modules.Planning == true {
		go database.DatabaseOpen()
		go database.CheckPlannedGames(&s)
	}
	if config.Cfg.Modules.Lottery == true {
		go PasswordLottery(s)
	}
	if config.Cfg.Modules.TimedChannelUnlock == true {
		go responder_functions.TimedChannelUnlock(s)
	}
	//Initialize the commands for Discord
	responder_functions.Ready(s, ready)

}

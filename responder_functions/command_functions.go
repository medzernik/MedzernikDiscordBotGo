// Package responder_functions This file contains the logic for the next-generation Commands, instead of the old prefix based responses.
package responder_functions

import (
	"bufio"
	"database/sql"
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"github.com/bwmarrin/discordgo"
	_ "github.com/dustin/go-humanize"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// ageJoinedCMD Checks the age of the user on join
func ageJoinedCMD(s *discordgo.Session, m *discordgo.InteractionCreate, cmd []interface{}) {

	//userId := m.ApplicationCommandData().Options[0].UserValue(s).ID
	userId := fmt.Sprintf("%s", cmd[0])

	//Every time a command is run, get a list of all users. This serves the purpose to then print the name of the corresponding user.
	//TODO: cache it in redis
	var membersCached []*discordgo.Member

	for i := range ReadyInfoPublic.Guilds {
		if ReadyInfoPublic.Guilds[i].ID == m.GuildID {
			membersCached = ReadyInfoPublic.Guilds[i].Members
		}
	}

	var userName string

	for i := range membersCached {
		if membersCached[i].User.ID == userId {
			userName = membersCached[i].User.Username
		} else if membersCached[i].User.ID != userId && membersCached[i].User.ID == "" {
			command.SendTextEmbedCommand(s, m.ChannelID, command.StatusBot.ERR, m.Data.Type().String()+" : not a number or a mention", discordgo.EmbedTypeRich)
			return
		}
	}

	userTimeRaw, err := SnowflakeTimestamp(userId)
	if err != nil {
		command.SendTextEmbedCommand(s, m.ChannelID, command.StatusBot.ERR, m.Data.Type().String()+" : not a number or a mention", discordgo.EmbedTypeRich)
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
	command.SendTextEmbedCommand(s, m.ChannelID, command.StatusBot.OK+userName, command.ParseStringToMentionID(userId)+" "+
		" has an account age of:\n"+
		""+rokyString+" rokov\n"+
		""+dnyString+" dni\n"+
		""+hodinyString+" hodin\n"+
		""+minutyString+" minut\n"+sekundyString+" sekund"+"<:peepoLove:687313976043765810>"+
		"", discordgo.EmbedTypeRich)
	return
}

func timeout(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	IDString := command.ParseMentionToString(m[0].(string))
	Until := time.Now().Add(time.Duration(m[1].(int64)) * time.Minute)

	fmt.Println(Until.String())

	err := s.GuildMemberTimeout(cmd.GuildID, IDString, &Until)
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error timeout-ing a user: "+err.Error(), discordgo.EmbedTypeRich)
		return
	}

	command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+" User timed out",
		"Timed out "+command.ParseStringToMentionID(m[0].(string))+" for "+strconv.FormatInt(m[1].(int64), 10)+" minutes"+"\n (Until: "+
			Until.Format(time.RFC822)+")",
		discordgo.EmbedTypeRich)
	return
}

func kickUserCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	var reasonExists bool
	var reason string
	authorisedAdmin, errPerm := command.MemberHasPermission(s, cmd.GuildID, cmd.Member.User.ID, discordgo.PermissionAdministrator)
	if errPerm != nil {
		fmt.Println(errPerm)
		return
	}

	if authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTH, "Error kicking a user - insufficient rights.", discordgo.EmbedTypeRich)
		return
	}

	var membersCached []*discordgo.Member

	for i := range ReadyInfoPublic.Guilds {
		if ReadyInfoPublic.Guilds[i].ID == cmd.GuildID {
			membersCached = ReadyInfoPublic.Guilds[i].Members
		}
	}
	if len(m) > 1 {
		reason = fmt.Sprintf("%s", m[1])
		reasonExists = true
	}

	var KickUserString string = command.ParseMentionToString(fmt.Sprintf("%s", m[0]))

	s.ChannelMessageSend(cmd.ChannelID, "**[PERM]** Permissions check complete.")

	if authorisedAdmin == true {
		for i := range membersCached {
			if membersCached[i].User.ID == KickUserString {
				if reasonExists == true {
					//DM the user of his kick + reason
					userNotifChanID, err0 := s.UserChannelCreate(KickUserString)
					if err0 != nil {
						s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error notifying the user of his kick")
					} else {
						s.ChannelMessageSend(userNotifChanID.ID, "You have been kicked from the server. Reason: "+reason)
					}

					//perform the kick itself
					err := s.GuildMemberDeleteWithReason(cmd.GuildID, KickUserString, reason)
					if err != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error kicking user ID "+membersCached[i].User.ID, discordgo.EmbedTypeRich)
						return
					}

					//log the kick
					command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"KICKED", "Kicked user "+membersCached[i].User.Username+"for "+
						""+reason, discordgo.EmbedTypeRich)
					s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "User "+KickUserString+" ,was kicked for: "+fmt.Sprintf("%s", m[1])+" .Kicked by "+cmd.Member.Nick)
				} else {
					//DM the user of his kick
					userNotifChanID, err0 := s.UserChannelCreate(KickUserString)
					if err0 != nil {
						s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error notifying the user of his kick")
					} else {
						s.ChannelMessageSend(userNotifChanID.ID, "You have been kicked from the server.")
					}

					//perform the kick itself
					err := s.GuildMemberDelete(cmd.GuildID, KickUserString)
					if err != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error kicking user ID "+membersCached[i].User.ID, discordgo.EmbedTypeRich)
						return
					}

					//log the kick
					//TODO: Fix the KickUserString -> stringf(m[0])
					command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"KICKED", "Kicked user "+membersCached[i].User.Username, discordgo.EmbedTypeRich)
					s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "User "+KickUserString+" "+fmt.Sprintf("%s", m[0])+" Kicked by "+cmd.Member.Nick)
				}
			}

		}
	}
	return
}

// banUserCMD bans a user
func banUserCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	var reason string
	var reasonExists bool = false
	var daysDelete uint64 = 7
	authorisedAdmin, errPerm := command.MemberHasPermission(s, cmd.GuildID, cmd.Member.User.ID, discordgo.PermissionAdministrator)
	if errPerm != nil {
		fmt.Println(errPerm)
		return
	}

	if authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTH, "Error banning a user - insufficient rights.", discordgo.EmbedTypeRich)
		return
	}

	var membersCached []*discordgo.Member

	for i := range ReadyInfoPublic.Guilds {
		if ReadyInfoPublic.Guilds[i].ID == cmd.GuildID {
			membersCached = ReadyInfoPublic.Guilds[i].Members
		}
	}

	var BanUserString string = fmt.Sprintf("%s", m[0])
	if len(m) > 1 {
		reason = fmt.Sprintf("%s", m[1])
		reasonExists = true
	}
	if len(m) > 2 {
		daysDelete = m[2].(uint64)
	}

	s.ChannelMessageSend(cmd.ChannelID, "**[PERM]** Permissions check complete.")

	if authorisedAdmin == true {
		for i := range membersCached {
			if membersCached[i].User.ID == BanUserString {
				if reasonExists == true {
					userNotifChanID, err0 := s.UserChannelCreate(BanUserString)
					if err0 != nil {
						s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error notifying the user of his ban")
					} else {
						s.ChannelMessageSend(userNotifChanID.ID, "You have been banned from the server. Reason: "+reason)
					}

					err := s.GuildBanCreateWithReason(cmd.GuildID, BanUserString, reason, int(daysDelete))
					if err != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error banning user ID "+membersCached[i].User.ID, discordgo.EmbedTypeRich)
						return
					}
					command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"BANNED", "Banned user "+membersCached[i].User.Username+"for "+
						""+reason, discordgo.EmbedTypeRich)
				} else {
					userNotifChanID, err0 := s.UserChannelCreate(BanUserString)
					if err0 != nil {
						s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[ERR]** Error notifying the user of his ban")
					} else {
						s.ChannelMessageSend(userNotifChanID.ID, "You have been banned from the server.")
					}

					err1 := s.GuildBanCreate(cmd.GuildID, BanUserString, int(daysDelete))
					if err1 != nil {
						command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error banning user ID "+membersCached[i].User.ID, discordgo.EmbedTypeRich)
						return
					}
					command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"BANNED", "Banning user "+membersCached[i].User.Username, discordgo.EmbedTypeRich)
					//TODO: Fix the BanUserString -> stringf(m[0])
					s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "User "+BanUserString+" "+fmt.Sprintf("%s", m[0])+" Banned by "+cmd.Member.Nick)
				}
			}

		}
	}
	return
}

// checkUsersCMD Checks the age of users
func checkUsersCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	var timeToCheckUsers int64
	if len(m) > 0 {
		timeToCheckUsers = m[0].(int64)
		timeToCheckUsers *= -1
	} else {
		timeToCheckUsers = 24 * -1
	}

	//variable definitions
	authorisedAdmin, errPerm := command.MemberHasPermission(s, cmd.GuildID, cmd.Member.User.ID, discordgo.PermissionAdministrator)
	if errPerm != nil {
		fmt.Println(errPerm)
		return
	}

	if authorisedAdmin == true {
		var membersCached []*discordgo.Member

		for i := range ReadyInfoPublic.Guilds {
			if ReadyInfoPublic.Guilds[i].ID == cmd.GuildID {
				membersCached = ReadyInfoPublic.Guilds[i].Members
			}
		}
		var mainOutputMsg string
		var IDOutputMsg string

		//iterate over the members_cached array. Maximum limit is 1000.
		for i := range membersCached {
			userTimeJoin := membersCached[i].JoinedAt
			var timeVar int64 = int64(userTimeJoin.Sub(time.Now()).Hours())

			if timeVar > timeToCheckUsers {
				mainOutputMsg += "This user is too young (less than " +
					"" + strconv.FormatInt(timeToCheckUsers*-1, 10) + "h join age): " +
					"" + membersCached[i].User.Username + " ,**ID:** " +
					"" + membersCached[i].User.ID + "\n"
				IDOutputMsg += membersCached[i].User.ID + " "
			}
		}
		//print out the amount of members_cached (max is currently 1000)
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"RECENT USERS", mainOutputMsg+"\n**IDs of the users (copyfriendly):**\n"+IDOutputMsg, discordgo.EmbedTypeRich)
	} else if authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTH, "You do not have the permission to use this command", discordgo.EmbedTypeRich)
		return
	}
	return

}

// planGameCMD Plans a game for a person with a timed reminder
func planGameCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	//go GamePlanInsertCMD(s, cmd, m)
	return
}

// plannedGamesCMD Checks the planned games and outputs them into the guild
func plannedGamesCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate) {
	//open database and then close it (defer)
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer func(sqliteDatabase *sql.DB) {
		err := sqliteDatabase.Close()
		if err != nil {
			fmt.Println("error closing the database: ", err)
		}
	}(sqliteDatabase)

	var plannedGames string
	//database.DisplayAllGamesPlanned(sqliteDatabase, &plannedGames)

	//send info to channel
	command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"PLANNED GAMES", plannedGames, discordgo.EmbedTypeRich)
	return
}

/*
// GamePlanInsertCMD Inserts the game into the database
func GamePlanInsertCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	//open database and then close it (defer)
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer func(sqliteDatabase *sql.DB) {
		err := sqliteDatabase.Close()
		if err != nil {

		}
	}(sqliteDatabase)

	//transform to timestamp
	splitTimeArgument := strings.Split(fmt.Sprintf("%s", m[0]), ":")

	//Put hours into timeHours
	timeHour, err := strconv.Atoi(splitTimeArgument[0])
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error converting hours", discordgo.EmbedTypeRich)
		return
	}
	//put minutes into timeMinute
	timeMinute, err := strconv.Atoi(splitTimeArgument[1])
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error converting minutes", discordgo.EmbedTypeRich)
		//(*s).ChannelMessageSend((*m).ChannelID, "**[ERR]** Error converting minutes")
		//fmt.Printf("%s", err)
		return
	}
	//get current date and replace hours and minutes with user variables
	gameTimestamp := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), timeHour, timeMinute, time.Now().Second(), 0, time.Now().Location())
	//gameTimestampInt := gameTimestamp.Unix()

	//export to database
	//database.InsertGame(sqliteDatabase, gameTimestampInt, fmt.Sprintf("%s", m[1]), fmt.Sprintf("%s", m[2]))

	var plannedGames string
	//database.DisplayGamePlanned(sqliteDatabase, &plannedGames)

	command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"PLANNED A GAME", plannedGames, discordgo.EmbedTypeRich)
	return
}

*/

// topicCMD Outputs a random topic for discussion found in topic_questions.txt
func topicCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate) {
	fileHandle, err := os.Open("topic_questions.txt")
	if err != nil {
		fmt.Println("error reading the file: ", err)
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error reading the file topic_questions.txt", discordgo.EmbedTypeRich)
		return
	}
	defer func(fileHandle *os.File) {
		err := fileHandle.Close()
		if err != nil {
			fmt.Println("**[ERR]** error closing the file with topics")
		}
	}(fileHandle)

	fileScanner := bufio.NewScanner(fileHandle)

	var splitTopic []string

	for fileScanner.Scan() {
		splitTopic = append(splitTopic, fileScanner.Text())
	}

	//a, b is the length of the topic.
	a := 0
	b := len(splitTopic)

	rand.Seed(time.Now().UnixNano())
	n := a + rand.Intn(b-a+1)

	//this checks the slice length and prevents a panic if for any chance it happened. Just in case.
	if n > len(splitTopic)-1 {
		fmt.Println("**[ERR_PARSE]** Slice is smaller than allowed\n This error should not have ever happened...")
		return
	}

	command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"TOPIC", splitTopic[n], discordgo.EmbedTypeRich)
	return
}

// getWeatherCMD outputs weather information from openWeatherMap
func getWeatherCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	type wData struct {
		name       string
		weather    string
		condition  string
		temp       string
		tempMax    string
		tempMin    string
		tempFeel   string
		pressure   string
		humidity   string
		windSpeed  string
		rainAmount string
		sunrise    string
		sunset     string
	}

	w, err := owm.NewCurrent("C", "en", config.Cfg.ServerInfo.WeatherAPIKey)
	if err != nil {
		fmt.Println("Error processing the request")
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error processing the request", discordgo.EmbedTypeRich)
	}

	var commandString string = fmt.Sprintf("%s", m[0])

	err2 := w.CurrentByName(commandString)
	if err2 != nil {
		log.Println(err2)
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "The city "+commandString+" does not exist", discordgo.EmbedTypeRich)
		return
	}

	var weatherData = wData{
		name:       w.Name,
		weather:    w.Weather[0].Main,
		condition:  w.Weather[0].Description,
		temp:       strconv.FormatFloat(w.Main.Temp, 'f', 1, 64) + " °C",
		tempMax:    strconv.FormatFloat(w.Main.TempMax, 'f', 1, 64) + " °C",
		tempMin:    strconv.FormatFloat(w.Main.TempMin, 'f', 1, 64) + " °C",
		tempFeel:   strconv.FormatFloat(w.Main.FeelsLike, 'f', 1, 64) + " °C",
		pressure:   strconv.FormatFloat(w.Main.Pressure, 'f', 1, 64) + " hPa",
		humidity:   strconv.FormatInt(int64(w.Main.Humidity), 10) + " %",
		windSpeed:  strconv.FormatFloat(w.Wind.Speed, 'f', 1, 64) + " km/h",
		rainAmount: strconv.FormatFloat(w.Rain.OneH*10, 'f', 1, 64) + " %",
		sunrise:    time.Unix(int64(w.Sys.Sunrise), 0).Format(time.Kitchen),
		sunset:     time.Unix(int64(w.Sys.Sunset), 0).Format(time.Kitchen),
	}

	embed := command.NewEmbed().
		SetTitle("WEATHER IN: "+strings.ToUpper(weatherData.name)).
		SetDescription(":cloud: **"+strings.ToUpper(weatherData.condition)+"**").
		AddField("**TEMPERATURE**",
			":thermometer: "+weatherData.temp+
				"\n:high_brightness: "+weatherData.tempMax+
				"\n:low_brightness: "+weatherData.tempMin).
		AddField("**PRESSURE**", ":diving_mask: "+weatherData.pressure).
		AddField("**HUMIDITY**", ":sweat: "+weatherData.humidity).
		AddField("**WIND SPEED**", ":dash: "+weatherData.windSpeed).
		AddField("**RAINFALL**", ":droplet: "+weatherData.rainAmount).
		AddField("**SUN RISE/SET**",
			":sunrise_over_mountains: "+weatherData.sunrise+
				"\n:city_sunset: "+weatherData.sunset).
		InlineAllFields().
		SetColor(3066993).MessageEmbed

	s.ChannelMessageSendEmbed(cmd.ChannelID, embed)
	return

}

// purgeMessagesCMD Purges messages 1-100 in the current channel
func purgeMessagesCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	authorisedAdmin, errPerm := command.MemberHasPermission(s, cmd.GuildID, cmd.Member.User.ID, discordgo.PermissionAdministrator)
	if errPerm != nil {
		fmt.Println(errPerm)
		return
	}

	if authorisedAdmin == true {
		var messageArrayToDelete []string

		numMessages := m[0].(int64)
		if numMessages > 99 || numMessages < 1 {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.SYNTAX, "The min-max of the number is 1-100", discordgo.EmbedTypeRich)
			return
		}

		messageArrayComplete, err1 := s.ChannelMessages(cmd.ChannelID, int(numMessages), cmd.ID, "", "")
		if err1 != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Cannot get the ID of messages", discordgo.EmbedTypeRich)
			return
		}

		for i := range messageArrayComplete {
			messageArrayToDelete = append(messageArrayToDelete, messageArrayComplete[i].ID)
		}

		err2 := s.ChannelMessagesBulkDelete(cmd.ChannelID, messageArrayToDelete)
		if err2 != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error deleting the requested messages...", discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"PURGED", "Purged "+strconv.FormatInt(int64(len(messageArrayToDelete)), 10)+" "+
			"messages", discordgo.EmbedTypeRich)

		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[LOG]** User "+cmd.Member.Nick+" deleted "+strconv.FormatInt(int64(len(messageArrayToDelete)), 10)+" messages in channel "+"<#"+cmd.ChannelID+">")

		return
	} else {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTH, "Insufficient permissions.", discordgo.EmbedTypeRich)
		return
	}

}

// purgeMessagesCMDMessage This function purges messages for the Application MessageCommand interface
func purgeMessagesCMDMessage(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	authorisedAdmin, errPerm := command.MemberHasPermission(s, cmd.GuildID, cmd.Member.User.ID, discordgo.PermissionAdministrator)
	if errPerm != nil {
		fmt.Println(errPerm)
		return
	}

	if authorisedAdmin == true {
		var messageArrayToDelete []string

		messageArrayComplete, err1 := s.ChannelMessages(cmd.ChannelID, 0, cmd.ID, m[0].(string), "")
		if err1 != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Cannot get the ID of messages", discordgo.EmbedTypeRich)
			return
		}

		for i := range messageArrayComplete {
			messageArrayToDelete = append(messageArrayToDelete, messageArrayComplete[i].ID)
		}

		err2 := s.ChannelMessagesBulkDelete(cmd.ChannelID, messageArrayToDelete)
		if err2 != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error deleting the requested messages...", discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"PURGED", "Purged "+strconv.FormatInt(int64(len(messageArrayToDelete)), 10)+" "+
			"messages", discordgo.EmbedTypeRich)

		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[LOG]** User "+cmd.Member.Nick+" deleted "+strconv.FormatInt(int64(len(messageArrayToDelete)), 10)+" messages in channel "+"<#"+cmd.ChannelID+">")

		return
	} else {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTH, "Insufficient permissions.", discordgo.EmbedTypeRich)
		return
	}

}

// purgeMessagesCMDMessageOnlyAuthor This function serves to purge messages of only the clicked author through the ApplicationMessageCommand interface
func purgeMessagesCMDMessageOnlyAuthor(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	authorisedAdmin, errPerm := command.MemberHasPermission(s, cmd.GuildID, cmd.Member.User.ID, discordgo.PermissionAdministrator)
	if errPerm != nil {
		fmt.Println(errPerm)
		return
	}

	if authorisedAdmin == true {
		var messageArrayToDelete []string

		messageArrayComplete, err1 := s.ChannelMessages(cmd.ChannelID, 0, cmd.ID, m[0].(string), "")
		if err1 != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Cannot get the ID of messages", discordgo.EmbedTypeRich)
			return
		}

		messageArrayCompleteLastIndex, err3 := s.ChannelMessages(cmd.ChannelID, 1, "", "", m[0].(string))
		if err3 != nil {
			fmt.Println(err3)
		}

		for i := range messageArrayComplete {
			if messageArrayComplete[i].Author.ID == messageArrayCompleteLastIndex[0].Author.ID {
				messageArrayToDelete = append(messageArrayToDelete, messageArrayComplete[i].ID)
			}
		}

		err2 := s.ChannelMessagesBulkDelete(cmd.ChannelID, messageArrayToDelete)
		if err2 != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error deleting the requested messages...", discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"PURGED", "Purged "+
			""+strconv.FormatInt(int64(len(messageArrayToDelete)), 10)+" messages from user: "+
			""+messageArrayCompleteLastIndex[0].Author.Mention(), discordgo.EmbedTypeRich)

		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[LOG]** User "+cmd.Member.Nick+" deleted "+strconv.FormatInt(int64(len(messageArrayToDelete)), 10)+" messages in channel "+"<#"+cmd.ChannelID+">")

		return
	} else {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTH, "Insufficient permissions.", discordgo.EmbedTypeRich)
		return
	}

}

func getMemberCount(cmd *discordgo.InteractionCreate) int {
	var memberListLength int

	for i := range ReadyInfoPublic.Guilds {
		if ReadyInfoPublic.Guilds[i].ID == cmd.GuildID {
			memberListLength = ReadyInfoPublic.Guilds[i].MemberCount
			return memberListLength
		}
	}

	fmt.Println("Error getting the member list length.")
	return 0
}

/*
func getMemberList(cmd *discordgo.InteractionCreate) []*discordgo.Member {
	var membersList []*discordgo.Member

	for i := range ReadyInfoPublic.Guilds {
		if ReadyInfoPublic.Guilds[i].ID == cmd.GuildID {
			membersList = ReadyInfoPublic.Guilds[i].Members
			return membersList
		}
	}

	fmt.Println("Error getting the memberlist.")
	return nil
}

*/

// membersCMD outputs the number of current members of the server. No returns
func membersCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate) {

	memberListLength := getMemberCount(cmd)

	command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+strconv.FormatUint(uint64(memberListLength), 10), ""+
		"There are "+strconv.FormatUint(uint64(memberListLength), 10)+" members on the server", discordgo.EmbedTypeRich)

	return
}

// pruneCountCMD outputs the number of users that could be pruned
func pruneCountCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	pruneDaysInt := m[0].(int64)

	if pruneDaysInt < 7 || pruneDaysInt > 30 {
		pruneDaysInt = 0
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.WARN, "Command is limited to range 7-30 for safety reasons", discordgo.EmbedTypeRich)
		return
	}

	pruneDaysCount, err := s.GuildPruneCount(cmd.GuildID, uint32(pruneDaysInt))
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error checking members to prune.", discordgo.EmbedTypeRich)
		return
	}

	command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+strconv.FormatUint(uint64(pruneDaysCount), 10), ""+
		"There are "+strconv.FormatUint(uint64(pruneDaysCount), 10)+" members to prune", discordgo.EmbedTypeRich)

	return
}

// pruneMembersCMD prunes members
func pruneMembersCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {
	authorisedAdmin, errPerm := command.MemberHasPermission(s, cmd.GuildID, cmd.Member.User.ID, discordgo.PermissionAdministrator)
	if errPerm != nil {
		fmt.Println(errPerm)
		return
	}

	if authorisedAdmin == true {
		//request prune number amount
		pruneDaysCountInt := m[0].(int64)

		var pruneDaysCountUInt = uint32(pruneDaysCountInt)

		if pruneDaysCountInt == 0 {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.SYNTAX, "Cannot prune time of 0 days. Allowed frame is 7-30", discordgo.EmbedTypeRich)
			s.ChannelMessageSend(cmd.ChannelID, "**[ERR]** Invalid days to prune (0)")
			return
		}

		//prunes the members and assigns the result of the pruned members count to a variable
		prunedMembersCount, err1 := s.GuildPrune(cmd.GuildID, pruneDaysCountUInt)
		if err1 != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error pruning members", discordgo.EmbedTypeRich)
		}

		//log output

		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"PRUNED", strconv.FormatInt(int64(prunedMembersCount), 10)+
			" members from the server", discordgo.EmbedTypeRich)
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[LOG]** User "+cmd.Member.Nick+
			" used a prune and kicked "+strconv.FormatInt(int64(prunedMembersCount), 10)+" members")
		return

		//permission output
	} else {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTH, "Insufficient permissions", discordgo.EmbedTypeRich)
		return
	}

}

// setRoleChannelPermCMD sets a channel permission using an int value
func setRoleChannelPermCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	authorisedAdmin, errPerm := command.MemberHasPermission(s, cmd.GuildID, cmd.Member.User.ID, discordgo.PermissionAdministrator)
	if errPerm != nil {
		fmt.Println(errPerm)
		return
	}

	//check if user is admin before using the command
	if authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTH, "You do not have the permission to change the channel permissions for a role.", discordgo.EmbedTypeRich)
		return
	}

	//get the role set
	permissionRole := command.ParseRoleMentionToString(fmt.Sprintf("%s", m[1]))
	//get the ID set
	permissionID := m[2].(int64)
	//get whether to allow or deny the permissions
	permissionAllow := m[0].(bool)

	// GET WHETHER TO SET PERMS FOR A ROLE OR A MEMBER

	if permissionAllow == true {
		err := s.ChannelPermissionSet(cmd.ChannelID, permissionRole, 0, permissionID, 0)
		if err != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error setting the permissions on the channel."+err.Error(), discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"PERMISSIONS ALLOWED", "Permissions "+strconv.FormatInt(permissionID, 10)+" successfully allowed"+
			"", discordgo.EmbedTypeRich)
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Admin "+cmd.Member.Nick+" allowed permissions "+strconv.FormatInt(permissionID, 10)+""+
			" for channel "+command.ParseStringToChannelID(cmd.ChannelID)+" to a role "+command.ParseStringToRoleMention(permissionRole))
		return

		//if not allowed then not allow
	} else if permissionAllow == false {
		err := s.ChannelPermissionSet(cmd.ChannelID, permissionRole, 0, 0, permissionID)
		if err != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error setting the permissions on the channel."+err.Error(), discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"PERMISSIONS DENIED", "Permissions "+strconv.FormatInt(permissionID, 10)+" successfully denied"+
			"", discordgo.EmbedTypeRich)
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Admin "+cmd.Member.Nick+" denied permissions "+strconv.FormatInt(permissionID, 10)+""+
			" for channel "+command.ParseStringToChannelID(cmd.ChannelID)+" to a role "+command.ParseStringToRoleMention(permissionRole))
		return

		//if there is an invalid syntax with the allow/deny argument then
	} else {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.SYNTAX, "Invalid syntax, .setpermission **<allow,deny>** "+
			"@roletoset INTPERMIDS", discordgo.EmbedTypeRich)
		return
	}

}

// setUserChannelPermCMD sets a channel permission using an int value
func setUserChannelPermCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	authorisedAdmin, errPerm := command.MemberHasPermission(s, cmd.GuildID, cmd.Member.User.ID, discordgo.PermissionAdministrator)
	if errPerm != nil {
		fmt.Println(errPerm)
		return
	}

	//check if user is admin before using the command
	if authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTH, "You do not have the permission to change the channel permissions for a role.", discordgo.EmbedTypeRich)
		return
	}

	//get the role to set
	permissionRole := command.ParseRoleMentionToString(fmt.Sprintf("%s", m[1]))

	permissionID := m[2].(int64)

	//get whether to allow or deny the permissions
	permissionAllow := m[0].(bool)
	if permissionAllow == true {
		err := s.ChannelPermissionSet(cmd.ChannelID, permissionRole, 1, permissionID, 0)
		if err != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error setting the permissions on the channel."+err.Error(), discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"PERMISSIONS ALLOWED", "Permissions "+strconv.FormatInt(permissionID, 10)+" successfully allowed"+
			"", discordgo.EmbedTypeRich)
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Admin "+cmd.ChannelID+" denied permissions "+strconv.FormatInt(permissionID, 10)+""+
			" for channel "+command.ParseStringToChannelID(cmd.ChannelID)+" to a user "+command.ParseStringToMentionID(permissionRole))
		return

		//if not allowed then not allow
	} else if permissionAllow == false {
		err := s.ChannelPermissionSet(cmd.ChannelID, permissionRole, 1, 0, permissionID)
		if err != nil {
			command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error setting the permissions on the channel."+err.Error(), discordgo.EmbedTypeRich)
			return
		}
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"PERMISSIONS DENIED", "Permissions "+strconv.FormatInt(permissionID, 10)+" successfully denied"+
			"", discordgo.EmbedTypeRich)
		s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Admin "+cmd.Member.Nick+" denied permissions "+strconv.FormatInt(permissionID, 10)+""+
			" for channel "+command.ParseStringToChannelID(cmd.ChannelID)+" to a user "+command.ParseStringToMentionID(permissionRole))
		return
	}
}

//redirectDiscussionCMD  sets a channel to a big slowmode for 10 minutes and then redirects the conversation elsewhere. When threads become available, sets the thread and more...
func redirectDiscussionCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	authorisedAdmin, errPerm := command.MemberHasPermission(s, cmd.GuildID, cmd.Member.User.ID, discordgo.PermissionAdministrator)
	if errPerm != nil {
		fmt.Println(errPerm)
		return
	}

	//check if user is admin before using the command
	if authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTH, "You do not have the permission to change the channel slowmode.", discordgo.EmbedTypeRich)
		return
	}

	originalChannelInfo, err := s.Channel(cmd.ChannelID)
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Cannot get info of the channel to modify, aborting.", discordgo.EmbedTypeRich)
		return
	}

	var channelIDString string
	var slowmodeChannelSet discordgo.ChannelEdit = discordgo.ChannelEdit{
		Name:                 originalChannelInfo.Name,
		Topic:                originalChannelInfo.Topic,
		NSFW:                 originalChannelInfo.NSFW,
		Position:             originalChannelInfo.Position,
		Bitrate:              originalChannelInfo.Bitrate,
		UserLimit:            originalChannelInfo.UserLimit,
		PermissionOverwrites: originalChannelInfo.PermissionOverwrites,
		ParentID:             originalChannelInfo.ParentID,
		RateLimitPerUser:     360,
	}

	channelIDString = command.ParseChannelToString(fmt.Sprintf("%v", m[0]))

	_, err1 := s.ChannelEditComplex(cmd.ChannelID, &slowmodeChannelSet)
	if err1 != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error setting the slowmode.", discordgo.EmbedTypeRich)
		return
	}
	command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"SLOWMODE SET FOR "+strconv.FormatInt(int64(slowmodeChannelSet.RateLimitPerUser), 10)+" SECONDS", "Continue discussion in "+
		""+command.ParseStringToChannelID(channelIDString), discordgo.EmbedTypeRich)

	return
}

// slowModeChannelCMD sets a channel to a desired slowmode. 0 is a bugged value, so at least sets to 1 second.
func slowModeChannelCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	authorisedAdmin, errPerm := command.MemberHasPermission(s, cmd.GuildID, cmd.Member.User.ID, discordgo.PermissionAdministrator)
	if errPerm != nil {
		fmt.Println(errPerm)
		return
	}

	//check if user is admin before using the command
	if authorisedAdmin == false {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTH, "You do not have the permission to change the channel slowmode.", discordgo.EmbedTypeRich)
		return
	}

	//parse the argument to input
	numOfSeconds := m[0].(uint64)

	//verify inputs
	if numOfSeconds > 21600 {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Seconds must be in the valid range (0-21600)", discordgo.EmbedTypeRich)
		return
	} else if numOfSeconds == 0 {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTOFIX, "Due to the bug in discord, setting the number to 1 (smallest possible wait time)", discordgo.EmbedTypeRich)
		numOfSeconds = 0
	}

	//get the original channel info
	originalChannelInfo, err := s.Channel(cmd.ChannelID)
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Cannot get info of the channel to modify, aborting.", discordgo.EmbedTypeRich)
		return
	}

	//get the current channel info and only modify the slowmode var
	var slowmodeChannelSet discordgo.ChannelEdit = discordgo.ChannelEdit{
		Name:                 originalChannelInfo.Name,
		Topic:                originalChannelInfo.Topic,
		NSFW:                 originalChannelInfo.NSFW,
		Position:             originalChannelInfo.Position,
		Bitrate:              originalChannelInfo.Bitrate,
		UserLimit:            originalChannelInfo.UserLimit,
		PermissionOverwrites: originalChannelInfo.PermissionOverwrites,
		ParentID:             originalChannelInfo.ParentID,
		RateLimitPerUser:     int(numOfSeconds),
	}

	//set the slowmode
	channelEdited, err1 := s.ChannelEditComplex(cmd.ChannelID, &slowmodeChannelSet)
	if err1 != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Error setting the slowmode.", discordgo.EmbedTypeRich)
		return
	}

	//send the confirmation message
	command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"SLOWMODE "+strconv.FormatUint(numOfSeconds, 10)+""+
		" SECONDS", "Set the channel slowmode to "+strconv.FormatUint(numOfSeconds, 10)+""+
		" seconds per message", discordgo.EmbedTypeRich)

	s.ChannelMessageSend(config.Cfg.ChannelLog.ChannelLogID, "**[OK]** Admin "+cmd.Member.Nick+" set a "+strconv.FormatInt(int64(channelEdited.RateLimitPerUser), 10)+""+
		" seconds slowmode in channel "+command.ParseStringToChannelID(channelEdited.ID))
	return
}

// changeVoiceChannelCurrentCMD I hate it when the bot creates a temp voice, and I have to change the name and all the time
func changeVoiceChannelCurrentCMD(s *discordgo.Session, cmd *discordgo.InteractionCreate, m []interface{}) {

	//find the user current voice channel
	currentUserInfo, err := FindUserVoiceState(s, command.ParseMentionToString(cmd.Member.Mention()))
	if err != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "You are not joined in any "+
			"voice channel. Or the bot can't see you. Either one.", discordgo.EmbedTypeRich)
		return
	}

	//get the original channel info
	originalChannelInfo, err1 := s.Channel(currentUserInfo.ChannelID)
	if err1 != nil {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Cannot get any info of the"+
			"channel you are joined in. Bot probably can't see into the channel.", discordgo.EmbedTypeRich)
		return
	}
	//default variable
	var bitrate int = 384000
	var errorCounter uint = 0

	name := fmt.Sprintf("%s", m[0])

	if len(m) > 1 {
		bitrate = int(m[1].(uint64))
		bitrate = bitrate * 1000
	}
	if bitrate < 8 || bitrate > 384000 {
		command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.SYNTAX, "Bitrate can be from 8 to 384 kbps", discordgo.EmbedTypeRich)
	}

	var voiceChannelSet discordgo.ChannelEdit = discordgo.ChannelEdit{
		Name:                 name,
		Topic:                originalChannelInfo.Topic,
		NSFW:                 originalChannelInfo.NSFW,
		Position:             originalChannelInfo.Position,
		Bitrate:              bitrate,
		UserLimit:            originalChannelInfo.UserLimit,
		PermissionOverwrites: originalChannelInfo.PermissionOverwrites,
		ParentID:             originalChannelInfo.ParentID,
		RateLimitPerUser:     originalChannelInfo.RateLimitPerUser,
	}

	for {
		_, err2 := s.ChannelEditComplex(currentUserInfo.ChannelID, &voiceChannelSet)
		if err2 != nil {
			switch errorCounter {
			case 0:
				command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTOFIX, "Trying lower bitrate... 284kbps", discordgo.EmbedTypeRich)
				bitrate = 284000
				errorCounter += 1

			case 1:
				command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTOFIX, "Trying lower bitrate... 128kbps", discordgo.EmbedTypeRich)
				bitrate = 128000
				errorCounter += 1

			case 2:
				command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.AUTOFIX, "Trying lower bitrate... 96kbps", discordgo.EmbedTypeRich)
				bitrate = 96000
				errorCounter += 1

			case 3:
				command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.ERR, "Some weird error occured", discordgo.EmbedTypeRich)
				return
			default:
				break
			}
		}
		break
	}
	command.SendTextEmbedCommand(s, cmd.ChannelID, command.StatusBot.OK+"CHANGED NAME AND BITRATE", "**Channel:** "+
		""+name+"\n**Bitrate:** "+strconv.FormatInt(int64(bitrate), 10), discordgo.EmbedTypeRich)
	return
}

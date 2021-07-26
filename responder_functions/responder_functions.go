// Package responder_functions contains all the logic and basic config commands for the responder commands.
package responder_functions

import (
	"bufio"
	"database/sql"
	"fmt"
	owm "github.com/briandowns/openweathermap"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"github.com/medzernik/SlovakiaDiscordBotGo/database"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// GuildIDNumber SETTINGS
//guildID. Change this to represent your server. ID of channel and server, String data type
const GuildIDNumber string = "513274646406365184"
const AdminChannel string = "837987736416813076"
const LogChannel string = "868194202012508190"
const TrustedChannel string = "751069355621744742"
const RoleMuteID string = "684159104901709825"

const apiKey string = "65bb37a9ac2af9128d6ceaf670043b39"

const Version string = "0.4.0"

type CommandStatus struct {
	OK     string
	ERR    string
	SYNTAX string
	WARN   string
	AUTH   string
}

// CommandStatusBot is a variable to pass to the messageEmbed to make an emoji
var CommandStatusBot CommandStatus = CommandStatus{
	OK:     ":white_check_mark: ",
	ERR:    ":bangbang: ERR",
	SYNTAX: ":question: SYNTAX",
	WARN:   ":warning: WARN",
	AUTH:   ":no_entry: AUTH",
}

func Zasielkovna(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	err := command.VerifyArguments(&cmd)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())

		return
	}
	s.ChannelMessageSend(m.ChannelID, "OVER 200% <a:medzernikShake:814055147583438848>")

}

// AgeJoined Checks the age of the user on join
func AgeJoined(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	err := command.VerifyArguments(&cmd, command.RegexArg{Expression: `^<@!(\d+)>$`, CaptureGroup: 1})
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	userId := cmd.Arguments[0]
	//Every time a command is run, get a list of all users. This serves the purpose to then print the name of the corresponding user.
	//TODO: cache it in redis
	membersCached := GetMemberListFromGuild(s, GuildIDNumber)

	var userName string

	for i := range membersCached {
		if membersCached[i].User.ID == userId {
			userName = membersCached[i].User.Username
			fmt.Println(userName)
		} else if membersCached[i].User.ID != userId && membersCached[i].User.ID == "" {
			s.ChannelMessageSend(m.ChannelID, "**[ERR] Bad user ID")

		}
	}

	fmt.Println(userId)

	userTimeRaw, err := SnowflakeTimestamp(userId)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "**[ERR] Bad user ID")
		return
	}

	userTime := time.Now().Sub(userTimeRaw)

	dny := int64(userTime.Hours() / 24)
	hodiny := int64(userTime.Hours()) - dny*24
	minuty := int64(userTime.Minutes()) - int64(userTime.Hours())*60
	sekundy := int64(userTime.Seconds()) - int64(userTime.Minutes())*60

	dnyString := strconv.FormatInt(dny, 10)
	hodinyString := strconv.FormatInt(hodiny, 10)
	minutyString := strconv.FormatInt(minuty, 10)
	sekundyString := strconv.FormatInt(sekundy, 10)

	s.ChannelMessageSend(m.ChannelID, "**[OK]**"+"**"+userName+"**"+" je tu s nami už:\n"+dnyString+" dni\n"+hodinyString+" hodin\n"+minutyString+" minut\n"+sekundyString+" sekund"+"<:peepoLove:687313976043765810>")

}

// Mute Muting function
func Mute(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	//CHANGE THIS: Enable to enable trustedUser Muting feature
	var trustedMutingEnabled bool = false

	//Variable initiation
	var authorisedAdmin bool = false
	var authorisedTrusted bool = false
	authorisedAdmin = command.VerifyAdmin(s, m, &authorisedAdmin, cmd)
	authorisedTrusted = command.VerifyTrusted(s, m, &authorisedTrusted, cmd)

	timeToCheckUsers := 24.0 * -1.0

	//Arguments checking
	if len(cmd.Arguments) < 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "**[SYNTAX]** Insufficient arguments provided (help: **.mute @user**)")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error processing request")
			return
		}
		return
	}

	//Verify, if user has any rights at all
	if authorisedAdmin == false && authorisedTrusted == false {
		s.ChannelMessageSend(m.ChannelID, "**[PERM]** Error muting a user - insufficient rights for any operation.")
		return
	}

	//Added only after the first check of rights, to prevent spamming of the requests
	membersCached := GetMemberListFromGuild(s, GuildIDNumber)
	var MuteUserString string = command.ParseMentionToString(cmd.Arguments[0])

	//Verify for the admin role before muting.
	if authorisedAdmin == true {
		for i := range membersCached {
			if membersCached[i].User.ID == MuteUserString {
				//Try to mute

				err := s.GuildMemberMute(GuildIDNumber, MuteUserString, true)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "**[INFO]** User not in VC, cannot mute")
				}
				err2 := s.GuildMemberRoleAdd(GuildIDNumber, MuteUserString, RoleMuteID)
				if err2 != nil {
					s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error muting adding the muted role.")
				}
				s.ChannelMessageSend(m.ChannelID, "**[OK]** Muted user "+MuteUserString)
				s.ChannelMessageSend(LogChannel, "**[LOG]** Administrator user "+"<@!"+m.Author.ID+">"+" Muted user: "+"<@!"+membersCached[i].User.ID+">")
				return
			}

		}
	}

	//If not, verify for the role of Trusted to try to mute
	if authorisedTrusted == true && authorisedAdmin == false && trustedMutingEnabled == true {
		for i := range membersCached {
			userTimeJoin, _ := membersCached[i].JoinedAt.Parse()
			timevar := userTimeJoin.Sub(time.Now()).Hours()
			if membersCached[i].User.ID == MuteUserString && timevar > timeToCheckUsers {
				//Error checking
				err := s.GuildMemberMute(GuildIDNumber, MuteUserString, true)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "**[INFO]** User not in VC, cannot mute")
				}
				err2 := s.GuildMemberRoleAdd(GuildIDNumber, MuteUserString, RoleMuteID)
				if err2 != nil {
					s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error muting adding the muted role.")
				}
				s.ChannelMessageSend(m.ChannelID, "**[OK]** Muted user younger than 24 hours "+MuteUserString)
				s.ChannelMessageSend(LogChannel, "**[LOG]** Trusted user "+"<@!"+m.Author.ID+">"+" Muted user: "+"<@!"+membersCached[i].User.ID+">")
				return
			} else if membersCached[i].User.ID == MuteUserString && timevar < timeToCheckUsers {
				s.ChannelMessageSend(m.ChannelID, "**[PERM]** Trusted users cannot mute anyone who joined for 24+ hours")
				return
			}

		}

	} else if trustedMutingEnabled == false {
		s.ChannelMessageSend(m.ChannelID, "**[OFF]** Muting by Trusted users is currently disabled.")
	} else {
		s.ChannelMessageSend(m.ChannelID, "**[ERR]** Undefined permission error.")
	}
	return
}

func KickUser(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	var reason string
	var authorisedAdmin bool = false
	authorisedAdmin = command.VerifyAdmin(s, m, &authorisedAdmin, cmd)

	if len(cmd.Arguments) < 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "**[SYNTAX]** Insufficient arguments provided (help: **.kick @user <optional_reason>**)")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error processing request")
			return
		}
		return
	}

	if len(cmd.Arguments) > 1 {
		reason = command.JoinArguments(cmd)
	}

	if authorisedAdmin == false {
		s.ChannelMessageSend(m.ChannelID, "**[PERM]** Error kicking a user - insufficient rights for operation.")
		return
	}

	membersCached := GetMemberListFromGuild(s, GuildIDNumber)

	var KickUserString string = command.ParseMentionToString(cmd.Arguments[0])

	s.ChannelMessageSend(m.ChannelID, "**[PERM]** Permissions check complete.")

	if authorisedAdmin == true {
		for i := range membersCached {
			if membersCached[i].User.ID == KickUserString {
				if len(reason) > 1 {
					//DM the user of his kick + reason
					userNotifChanID, err0 := s.UserChannelCreate(KickUserString)
					if err0 != nil {
						s.ChannelMessageSend(LogChannel, "**[ERR]** Error notifying the user of his kick")
					}
					s.ChannelMessageSend(userNotifChanID.ID, "You have been kicked from the server. Reason: "+reason)

					//perform the kick itself
					err := s.GuildMemberDeleteWithReason(GuildIDNumber, KickUserString, reason)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error kicking user")
						return
					}

					//log the kick
					s.ChannelMessageSend(m.ChannelID, "**[OK]** Kicking user: "+KickUserString+". \nReason: "+reason)
					s.ChannelMessageSend(LogChannel, "User "+KickUserString+" "+cmd.Arguments[0]+" Kicked by "+m.Author.Username)
				} else {
					//DM the user of his kick
					userNotifChanID, err0 := s.UserChannelCreate(KickUserString)
					if err0 != nil {
						s.ChannelMessageSend(LogChannel, "**[ERR]** Error notifying the user of his kick")
					}
					s.ChannelMessageSend(userNotifChanID.ID, "You have been kicked from the server.")

					//perform the kick itself
					err := s.GuildMemberDelete(GuildIDNumber, KickUserString)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error kicking user")
						return
					}

					//log the kick
					s.ChannelMessageSend(m.ChannelID, "**[OK]** Kicking user: "+KickUserString+". \nReason: "+reason)
					s.ChannelMessageSend(LogChannel, "User "+KickUserString+" "+cmd.Arguments[0]+" Kicked by "+m.Author.Username)
				}
			}

		}
	}
	return
}

func BanUser(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	var reason string
	var daysDelete int = 7
	var authorisedAdmin bool = false
	authorisedAdmin = command.VerifyAdmin(s, m, &authorisedAdmin, cmd)

	if len(cmd.Arguments) < 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "**[SYNTAX]** Insufficient arguments provided (help: **.ban @user <optional reason>**)")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error processing request")
			return
		}
		return
	}

	if len(cmd.Arguments) > 2 {
		reason = command.JoinArguments(cmd)
	}

	if authorisedAdmin == false {
		s.ChannelMessageSend(m.ChannelID, "**[PERM]** Error Banning user - insufficient rights for operation.")
		return
	}

	membersCached := GetMemberListFromGuild(s, GuildIDNumber)

	var BanUserString string = command.ParseMentionToString(cmd.Arguments[0])

	s.ChannelMessageSend(m.ChannelID, "**[PERM]** Permissions check complete.")

	if authorisedAdmin == true {
		for i := range membersCached {
			if membersCached[i].User.ID == BanUserString {
				if len(reason) > 0 {
					userNotifChanID, err0 := s.UserChannelCreate(BanUserString)
					if err0 != nil {
						s.ChannelMessageSend(LogChannel, "**[ERR]** Error notifying the user of his ban")
					}
					s.ChannelMessageSend(userNotifChanID.ID, "You have been banned from the server. Reason: "+reason)

					err := s.GuildBanCreateWithReason(GuildIDNumber, BanUserString, reason, daysDelete)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error banning user")
						return
					}
					s.ChannelMessageSend(LogChannel, "User "+BanUserString+" "+cmd.Arguments[0]+" Banned by "+m.Author.Username)
					s.ChannelMessageSend(m.ChannelID, "**[OK]** Banning user: "+BanUserString+". \nReason: "+reason)
				} else {
					userNotifChanID, err0 := s.UserChannelCreate(BanUserString)
					if err0 != nil {
						s.ChannelMessageSend(LogChannel, "**[ERR]** Error notifying the user of his ban")
					}
					s.ChannelMessageSend(userNotifChanID.ID, "You have been banned from the server.")

					err1 := s.GuildBanCreate(GuildIDNumber, BanUserString, daysDelete)
					if err1 != nil {
						s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error banning user")
						return
					}
					s.ChannelMessageSend(LogChannel, "User "+BanUserString+" "+cmd.Arguments[0]+" Banned by "+m.Author.Username)
					s.ChannelMessageSend(m.ChannelID, "**[OK]** Banning user: "+BanUserString+". \nReason: "+reason)
				}
			}

		}
	}
	return
}

// CheckUsers Checks the age of users
func CheckUsers(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	err := command.VerifyArguments(&cmd)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	//variable definitions
	var authorisedAdmin bool = false
	authorisedAdmin = command.VerifyAdmin(s, m, &authorisedAdmin, cmd)

	if authorisedAdmin == true {
		membersCached := GetMemberListFromGuild(s, GuildIDNumber)
		var mainOutputMsg string
		var IDOutputMsg string
		timeToCheckUsers := 24.0 * -1.0

		//iterate over the members_cached array. Maximum limit is 1000.
		for i := range membersCached {
			userTimeJoin, _ := membersCached[i].JoinedAt.Parse()
			var timeVar float64 = userTimeJoin.Sub(time.Now()).Hours()

			if timeVar > timeToCheckUsers {
				mainOutputMsg += "This user is too young (less than 24h join age): " + membersCached[i].User.Username + " ,**ID:** " + membersCached[i].User.ID + "\n"
				IDOutputMsg += membersCached[i].User.ID + " "
			}
		}
		//print out the amount of members_cached (max is currently 1000)
		s.ChannelMessageSend(m.ChannelID, mainOutputMsg+"\n**IDs of the users (copyfriendly):**\n"+IDOutputMsg)
	} else if authorisedAdmin == false {
		s.ChannelMessageSend(m.ChannelID, "**[PERM]** You do not have the permissions to use this command.")
	}
	return

}

// PlanGame Plans a game for a person with a timed reminder
func PlanGame(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	if len(cmd.Arguments) < 3 {
		s.ChannelMessageSend(m.ChannelID, "**[SYNTAX]** Insufficient arguments. Provided "+strconv.FormatInt(int64(len(cmd.Arguments)), 10)+" , Expected at least 3")
		return
	}
	GamePlanInsert(&cmd, &s, &m)
	return
}

func PlannedGames(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	if len(cmd.Arguments) > 0 {
		s.ChannelMessageSend(m.ChannelID, "**[SYNTAX]** Insufficient arguments. Provided "+strconv.FormatInt(int64(len(cmd.Arguments)), 10)+" , Expected no arguments")
		return
	}
	//open database and then close it (defer)
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer func(sqliteDatabase *sql.DB) {
		err := sqliteDatabase.Close()
		if err != nil {
			fmt.Println("error closing the database: ", err)
		}
	}(sqliteDatabase)

	var plannedGames string
	database.DisplayAllGamesPlanned(sqliteDatabase, &plannedGames)

	//send info to channel
	(*s).ChannelMessageSend((*m).ChannelID, plannedGames)
	return
}

// GamePlanInsert Inserts the game into the database
func GamePlanInsert(c *command.Command, s **discordgo.Session, m **discordgo.MessageCreate) {
	//open database and then close it (defer)
	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
	defer func(sqliteDatabase *sql.DB) {
		err := sqliteDatabase.Close()
		if err != nil {

		}
	}(sqliteDatabase)

	//transform to timestamp
	splitTimeArgument := strings.Split(c.Arguments[0], ":")

	//TODO: Check the capacity if it's sufficient, otherwise the program is panicking every time...
	if cap(splitTimeArgument) < 1 {
		(*s).ChannelMessageSend((*m).ChannelID, "**[ERR]** Error parsing time")
		return
	}

	//Put hours into timeHours
	timeHour, err := strconv.Atoi(splitTimeArgument[0])
	if err != nil {
		(*s).ChannelMessageSend((*m).ChannelID, "**[ERR]** Error converting hours")
		fmt.Printf("%s", err)
		return
	}
	//put minutes into timeMinute
	timeMinute, err := strconv.Atoi(splitTimeArgument[1])
	if err != nil {
		(*s).ChannelMessageSend((*m).ChannelID, "**[ERR]** Error converting minutes")
		fmt.Printf("%s", err)
		return
	}
	//get current date and replace hours and minutes with user variables
	gameTimestamp := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), timeHour, timeMinute, time.Now().Second(), 0, time.Now().Location())
	gameTimestampInt := gameTimestamp.Unix()

	fmt.Println(gameTimestampInt)

	//export to database
	database.InsertGame(sqliteDatabase, gameTimestampInt, c.Arguments[1], c.Arguments[2])

	var plannedgames string
	database.DisplayGamePlanned(sqliteDatabase, &plannedgames)

	(*s).ChannelMessageSend((*m).ChannelID, plannedgames)
	return
}

// SnowflakeTimestamp Function to check the user's join date
func SnowflakeTimestamp(ID string) (t time.Time, err error) {
	i, err := strconv.ParseInt(ID, 10, 64)
	if err != nil {
		return
	}
	timestamp := (i >> 22) + 1420070400000
	t = time.Unix(0, timestamp*1000000)
	return
}

// GetMemberListFromGuild Gets the member info
func GetMemberListFromGuild(s *discordgo.Session, guildID string) []*discordgo.Member {
	membersList, err := s.GuildMembers(guildID, "0", 1000)
	if err != nil {
		fmt.Println("**[CONF_ERR]** Error getting the memberlist (probably invalid guildID): " + guildID)
	}

	return membersList

}

// CheckRegularSpamAttack Checks the server for spam attacks
func CheckRegularSpamAttack(s *discordgo.Session) {
	//variable definitons
	var membersCached = GetMemberListFromGuild(s, GuildIDNumber)
	var tempMsg string
	var spamCounter int64
	var checkInterval time.Duration = 90
	var timeToCheckUsers = 10 * -1.0

	for {
		//iterate over the members_cached array. Maximum limit is 1000.
		for i := range membersCached {
			userTimeJoin, _ := membersCached[i].JoinedAt.Parse()
			timeVar := userTimeJoin.Sub(time.Now()).Minutes()

			if timeVar > timeToCheckUsers {
				tempMsg += "**[ALERT]** RAID PROTECTION ALERT!: User" + membersCached[i].User.Username + "join age: " + strconv.FormatFloat(timeToCheckUsers, 'f', 0, 64) + "\n"
				spamCounter += 1
			}

		}
		if spamCounter > 4 {
			s.ChannelMessageSend(AdminChannel, "**[WARN]** Possible RAID ATTACK detected!!! (<@&513275201375698954>) ("+strconv.FormatInt(spamCounter, 10)+" users joined in the last "+strconv.FormatFloat(timeToCheckUsers, 'f', 0, 64)+" hours)")
		}
		spamCounter = 0
		time.Sleep(checkInterval * time.Second)
	}

}

func Topic(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	err := command.VerifyArguments(&cmd)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	fileHandle, err := os.Open("topic_questions.txt")
	if err != nil {
		fmt.Println("error reading the file: ", err)
		s.ChannelMessageSend(m.ChannelID, err.Error())
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
		fmt.Println("[ERR_PARSE] Slice is smaller than allowed\n This error should not have ever happened...")
		return
	}

	s.ChannelMessageSend(m.ChannelID, splitTopic[n])
	return
}

func Fox(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "<a:medzernikShake:814055147583438848>")
}

// GetWeather outputs weather information from openWeatherMap
func GetWeather(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {

	if len(cmd.Arguments) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Usage:\n``.weather CityName``")
		return
	}

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

	w, err := owm.NewCurrent("C", "en", apiKey)
	if err != nil {
		fmt.Println("Error processing the request")
		log.Fatalln(err)
	}

	var commandString string = command.JoinArguments(cmd)

	err2 := w.CurrentByName(commandString)
	if err2 != nil {
		log.Println(err2)
		s.ChannelMessageSend(m.ChannelID, "**error: the city **"+commandString+"**does not exist.**")
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

	var weatherDataString string = "```\n" +
		"City:\t\t" + weatherData.name + "\n" +
		"Weather:\t" + weatherData.weather + "\n" +
		"Condition:\t" + weatherData.condition + "\n" +
		"Temperature:" + weatherData.temp + "\n" +
		"Max Temp:\t" + weatherData.tempMax + "\n" +
		"Min Temp:\t" + weatherData.tempMin + "\n" +
		"Feel Temp:\t" + weatherData.tempFeel + "\n" +
		"Pressure:\t" + weatherData.pressure + "\n" +
		"Humidity:\t" + weatherData.humidity + "\n" +
		"Wind Speed:\t" + weatherData.windSpeed + "\n" +
		"Rainfall:\t" + weatherData.rainAmount + "\n" +
		"Sunrise:\t" + weatherData.sunrise + "\n" +
		"Sunset:\t" + weatherData.sunset + "\n" +
		"```"

	s.ChannelMessageSend(m.ChannelID, weatherDataString)
	return

}

// TimedChannelUnlock automatically locks and unlocks a trusted channel
func TimedChannelUnlock(s *discordgo.Session) {
	var checkInterval time.Duration = 60

	fmt.Println("[INIT OK] Channel unlock system module initialized")

	for {
		if time.Now().Weekday() == time.Friday && time.Now().Hour() == 18 && time.Now().Minute() == 00 {
			//Unlock the channel
			//TargetType 0 = roleID, 1 = memberID
			err1 := s.ChannelPermissionSet(TrustedChannel, command.AuthorisedIDTrusted1, 0, 2251673408, 0)
			if err1 != nil {
				s.ChannelMessageSend(LogChannel, "**[ERR]** Error changing the permissions for role "+"<@"+command.AuthorisedIDTrusted1+">")
			}
			err2 := s.ChannelPermissionSet(TrustedChannel, command.AuthorisedIDTrusted2, 0, 2251673408, 0)
			if err2 != nil {
				s.ChannelMessageSend(LogChannel, "**[ERR]** Error changing the permissions for role "+"<@"+command.AuthorisedIDTrusted2+">")
			}
			err3 := s.ChannelPermissionSet(TrustedChannel, command.AuthorisedIDTrusted3, 0, 2251673408, 0)
			if err3 != nil {
				s.ChannelMessageSend(LogChannel, "**[ERR]** Error changing the permissions for role "+"<@"+command.AuthorisedIDTrusted3+">")
			}
			err4 := s.ChannelPermissionSet(TrustedChannel, command.AuthorisedIDTrustedMini, 0, 2251673408, 0)
			if err4 != nil {
				s.ChannelMessageSend(LogChannel, "**[ERR]** Error changing the permissions for role "+"<@"+command.AuthorisedIDTrustedMini+">")
			}

			fmt.Println("[OK] Opened the channel " + TrustedChannel)
		} else if time.Now().Weekday() == time.Monday && time.Now().Hour() == 6 && time.Now().Minute() == 0 {
			//Lock the channel
			//TargetType 0 = roleID, 1 = memberID
			err1 := s.ChannelPermissionSet(TrustedChannel, command.AuthorisedIDTrusted1, 0, 0, 2251673408)
			if err1 != nil {
				s.ChannelMessageSend(LogChannel, "**[ERR]** Error changing the permissions for role "+"<@"+command.AuthorisedIDTrusted1+">")
			}
			err2 := s.ChannelPermissionSet(TrustedChannel, command.AuthorisedIDTrusted2, 0, 0, 2251673408)
			if err2 != nil {
				s.ChannelMessageSend(LogChannel, "**[ERR]** Error changing the permissions for role "+"<@"+command.AuthorisedIDTrusted2+">")
			}
			err3 := s.ChannelPermissionSet(TrustedChannel, command.AuthorisedIDTrusted3, 0, 0, 2251673408)
			if err3 != nil {
				s.ChannelMessageSend(LogChannel, "**[ERR]** Error changing the permissions for role "+"<@"+command.AuthorisedIDTrusted3+">")
			}
			err4 := s.ChannelPermissionSet(TrustedChannel, command.AuthorisedIDTrustedMini, 0, 0, 2251673408)
			if err4 != nil {
				s.ChannelMessageSend(LogChannel, "**[ERR]** Error changing the permissions for role "+"<@"+command.AuthorisedIDTrustedMini+">")
			}
			fmt.Println("[OK] Closed the channel " + TrustedChannel)
		}

		time.Sleep(checkInterval * time.Second)
	}

}

func PurgeMessages(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {

	if len(cmd.Arguments) < 1 {
		s.ChannelMessageSend(m.ChannelID, "Usage: **.purge number**")
		return
	}

	var authorisedAdmin bool = false
	authorisedAdmin = command.VerifyAdmin(s, m, &authorisedAdmin, cmd)

	if authorisedAdmin == true {
		var messageArrayToDelete []string

		numMessages, err1 := strconv.ParseInt(cmd.Arguments[0], 10, 64)
		if err1 != nil {
			s.ChannelMessageSend(m.ChannelID, "**[ERR]** Invalid number provided")
			return
		}
		if numMessages > 99 || numMessages < 1 {
			s.ChannelMessageSend(m.ChannelID, "**[SYNTAX]** The number of messages must be 1-100")
			return
		}

		messageArrayComplete, err1 := s.ChannelMessages(m.ChannelID, int(numMessages), m.ID, "", "")
		if err1 != nil {
			s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error getting the ID of messages")
			return
		}

		for i := range messageArrayComplete {
			messageArrayToDelete = append(messageArrayToDelete, messageArrayComplete[i].ID)
		}

		err2 := s.ChannelMessagesBulkDelete(m.ChannelID, messageArrayToDelete)
		if err2 != nil {
			s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error deleting the requested messages...")
			return
		}
		s.ChannelMessageSend(m.ChannelID, "**[OK]** Deleted "+strconv.FormatInt(int64(len(messageArrayToDelete)), 10)+" messages")

		s.ChannelMessageSend(LogChannel, "**[LOG]** User "+"<@!"+m.Author.ID+">"+" deleted "+strconv.FormatInt(int64(len(messageArrayToDelete)), 10)+" messages in channel "+"<#"+m.ChannelID+">")

		//command.SendTextEmbed(s, m, "[OK]")

		return
	} else {
		s.ChannelMessageSend(m.ChannelID, "**[PERM]** Insufficient permissions.")
		return
	}

}

// Members outputs the number of current members of the server
func Members(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) uint64 {
	if len(cmd.Arguments) > 0 {
		command.SendTextEmbed(s, m, CommandStatusBot.SYNTAX, "Usage: **.count**\n Automatically discarding arguments...", discordgo.EmbedTypeRich)
	}

	memberList := GetMemberListFromGuild(s, GuildIDNumber)
	memberListLength := uint64(len(memberList))

	return memberListLength
}

// PruneCount outputs the number of users that could be pruned
func PruneCount(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) uint32 {

	if len(cmd.Arguments) < 1 {
		s.ChannelMessageSend(m.ChannelID, "**[SYNTAX]** Usage: **.prunecount days**")
		return 0
	}

	pruneDaysString := cmd.Arguments[0]
	pruneDaysInt, err1 := strconv.ParseInt(pruneDaysString, 10, 64)
	if err1 != nil {
		s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error parsing the days argument as a number (uint32)")
		return 0
	}

	if pruneDaysInt < 7 {
		pruneDaysInt = 0
		s.ChannelMessageSend(m.ChannelID, "**[ERR]** Must be at least 7 days of inactivity.")
		return 0
	}

	pruneDaysCount, err2 := s.GuildPruneCount(GuildIDNumber, uint32(pruneDaysInt))
	if err2 != nil {
		return 0
	}
	return pruneDaysCount

}

// PruneMembers prunes members
func PruneMembers(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	if len(cmd.Arguments) < 100 {
		s.ChannelMessageSend(m.ChannelID, "**[SYNTAX]** Usage: **.prunemembers days**")
		return
	}

	var authorisedAdmin bool = false
	authorisedAdmin = command.VerifyAdmin(s, m, &authorisedAdmin, cmd)

	if authorisedAdmin == true {
		//request prune number amount
		pruneDaysCountInt, err0 := strconv.ParseInt(cmd.Arguments[0], 10, 32)

		if err0 != nil {
			s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error parsing the number of days ")
		}

		var pruneDaysCountUInt = uint32(pruneDaysCountInt)

		if pruneDaysCountInt == 0 {
			s.ChannelMessageSend(m.ChannelID, "**[ERR]** Invalid days to prune (0)")
			return
		}

		//prunes the members and assigns the result of the pruned members count to a variable
		prunedMembersCount, err1 := s.GuildPrune(GuildIDNumber, pruneDaysCountUInt)
		if err1 != nil {
			s.ChannelMessageSend(m.ChannelID, "**[ERR]** Error pruning members")
		}

		//log output
		s.ChannelMessageSend(LogChannel, "**[LOG]** User "+"<@!"+m.Author.ID+">"+
			" used a prune and kicked "+strconv.FormatInt(int64(prunedMembersCount), 10)+" members")
		s.ChannelMessageSend(m.ChannelID, "**[OK]** pruned "+strconv.FormatInt(int64(prunedMembersCount), 10)+
			" members from the server")
		return

		//permission output
	} else {
		s.ChannelMessageSend(m.ChannelID, "**[PERM]** Insufficient permissions for operation")
		return
	}

}

// MassKick mass kicks user IDs
func MassKick(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {

}

// MassBan mass bans user IDs
func MassBan(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {

}

// SetChannelPermission sets a channel permission using an int value
func SetChannelPermission(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {

}

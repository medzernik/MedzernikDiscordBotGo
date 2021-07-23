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
const LogChannel string = "513280604507340804"

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

	for itera := range membersCached {
		if membersCached[itera].User.ID == userId {
			userName = membersCached[itera].User.Username
			fmt.Println(userName)
		} else if membersCached[itera].User.ID != userId && membersCached[itera].User.ID == "" {
			s.ChannelMessageSend(m.ChannelID, "Zlé ID slováka")

		}
	}

	fmt.Println(userId)

	userTimeRaw, err := SnowflakeTimestamp(userId)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Zlé ID slováka")
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

	s.ChannelMessageSend(m.ChannelID, "**"+userName+"**"+" je tu s nami už:\n"+dnyString+" dni\n"+hodinyString+" hodin\n"+minutyString+" minut\n"+sekundyString+" sekund"+"<:peepoLove:687313976043765810>")

}

// Mute Muting function
func Mute(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate, err error) {
	//enable to enable trustedUser Muting feature
	var trustedMutingEnabled bool = false

	timeToCheckUsers := 24.0 * -1.0
	var roleMuteID = "684159104901709825"

	//Arguments checking
	if len(cmd.Arguments) < 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "[SYNTAX] Insufficient arguments provided (help: **.mute @user**)")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "[ERR] Error processing request")
			return
		}
		return
	}

	//Verify, if user has any rights at all
	if command.VerifyAdmin(s, m) == false && command.VerifyTrusted(s, m) == false {
		s.ChannelMessageSend(m.ChannelID, "[PERM] Error muting a user - insufficient rights for any operation.")
		return
	}

	//Added only after the first check of rights, to prevent spamming of the requests
	membersCached := GetMemberListFromGuild(s, GuildIDNumber)
	var MuteUserString string = command.ParseMentionToString(cmd.Arguments[0])

	//Verify for the admin role before muting.
	if command.VerifyAdmin(s, m) == true {
		for i := range membersCached {
			if membersCached[i].User.ID == MuteUserString {
				s.ChannelMessageSend(m.ChannelID, "[OK] Muted user "+MuteUserString)
				err := s.GuildMemberMute(GuildIDNumber, MuteUserString, true)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "[ERR] Error muting the user in voice chat.")
				}
				err2 := s.GuildMemberRoleAdd(GuildIDNumber, MuteUserString, roleMuteID)
				if err2 != nil {
					s.ChannelMessageSend(m.ChannelID, "[ERR] Error muting adding the muted role.")
				}
				s.ChannelMessageSend(LogChannel, "[OK] Administrator user"+m.Author.Username+" Muted user: "+membersCached[i].Nick)
				return
			}
		}
	}

	//If not, verify for the role of Trusted to try to mute
	if command.VerifyTrusted(s, m) == true && command.VerifyAdmin(s, m) == false && trustedMutingEnabled == true {
		for i := range membersCached {
			userTimeJoin, _ := membersCached[i].JoinedAt.Parse()
			timevar := userTimeJoin.Sub(time.Now()).Hours()
			if membersCached[i].User.ID == MuteUserString && timevar > timeToCheckUsers {
				s.ChannelMessageSend(m.ChannelID, "[OK] Muted user younger than 24 hours "+MuteUserString)
				err := s.GuildMemberMute(GuildIDNumber, MuteUserString, true)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "[ERR] Error muting the user in voice chat.")
				}
				err2 := s.GuildMemberRoleAdd(GuildIDNumber, MuteUserString, roleMuteID)
				if err2 != nil {
					s.ChannelMessageSend(m.ChannelID, "[ERR] Error muting adding the muted role.")
				}
				s.ChannelMessageSend(LogChannel, "[OK] Trusted user"+m.Author.Username+" Muted user: "+membersCached[i].Nick)
				return
			} else if membersCached[i].User.ID == MuteUserString && timevar < timeToCheckUsers {
				s.ChannelMessageSend(m.ChannelID, "[PERM] Trusted users cannot mute anyone who joined for 24+ hours")
				return
			}

		}

	} else {
		s.ChannelMessageSend(m.ChannelID, "[OFF] Muting by Trusted users is currently disabled.")
	}
	return
}

func KickUser(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	var reason string

	if len(cmd.Arguments) < 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "[SYNTAX] Insufficient arguments provided (help: **.kick @user <optional_reason>**)")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "[ERR] Error processing request")
			return
		}
		return
	}

	if len(cmd.Arguments) > 1 {
		reason = command.JoinArguments(cmd)
	}

	if command.VerifyAdmin(s, m) == false {
		s.ChannelMessageSend(m.ChannelID, "[PERM] Error kicking a user - insufficient rights for operation.")
		return
	}

	membersCached := GetMemberListFromGuild(s, GuildIDNumber)

	var KickUserString string = command.ParseMentionToString(cmd.Arguments[0])

	s.ChannelMessageSend(m.ChannelID, "[PERM] Permissions check complete.")

	if command.VerifyAdmin(s, m) == true {
		for i := range membersCached {
			if membersCached[i].User.ID == KickUserString {
				if len(reason) > 1 {
					err := s.GuildMemberDeleteWithReason(GuildIDNumber, KickUserString, reason)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "[ERR] Error kicking user")
						return
					}
					s.ChannelMessageSend(m.ChannelID, "[OK] Kicking user: "+KickUserString+". \nReason: "+reason)
					s.ChannelMessageSend(LogChannel, "User "+KickUserString+" "+cmd.Arguments[0]+" Kicked by "+m.Author.Username)
				} else {
					err := s.GuildMemberDelete(GuildIDNumber, KickUserString)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "[ERR] Error kicking user")
						return
					}
					s.ChannelMessageSend(m.ChannelID, "[OK] Kicking user: "+KickUserString+". \nReason: "+reason)
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

	if len(cmd.Arguments) < 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "[SYNTAX] Insufficient arguments provided (help: **.kick @user <optional_reason>**)")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "[ERR] Error processing request")
			return
		}
		return
	}

	if len(cmd.Arguments) > 2 {
		reason = command.JoinArguments(cmd)
		//daysDelete, _ = strconv.Atoi(cmd.Arguments[2])
	}

	if command.VerifyAdmin(s, m) == false {
		s.ChannelMessageSend(m.ChannelID, "[PERM] Error Banning user - insufficient rights for operation.")
		return
	}

	membersCached := GetMemberListFromGuild(s, GuildIDNumber)

	var BanUserString string = command.ParseMentionToString(cmd.Arguments[0])

	s.ChannelMessageSend(m.ChannelID, "[PERM] Permissions check complete.")

	if command.VerifyAdmin(s, m) == true {
		for i := range membersCached {
			if membersCached[i].User.ID == BanUserString {
				if len(reason) > 0 {
					err := s.GuildBanCreateWithReason(GuildIDNumber, BanUserString, reason, daysDelete)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "[ERR] Error banning user")
						return
					}
					s.ChannelMessageSend(LogChannel, "User "+BanUserString+" "+cmd.Arguments[0]+" Banned by "+m.Author.Username)
					s.ChannelMessageSend(m.ChannelID, "[OK] Banning user: "+BanUserString+". \nReason: "+reason)
				} else {
					err := s.GuildBanCreate(GuildIDNumber, BanUserString, daysDelete)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, "[ERR] Error banning user")
						return
					}
					s.ChannelMessageSend(LogChannel, "User "+BanUserString+" "+cmd.Arguments[0]+" Banned by "+m.Author.Username)
					s.ChannelMessageSend(m.ChannelID, "[OK] Banning user: "+BanUserString+". \nReason: "+reason)
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
	//variable definitons
	//members_cached, _ := s.GuildMembers("513274646406365184", "0", 1000)
	membersCached := GetMemberListFromGuild(s, GuildIDNumber)
	var tempMsg string
	timeToCheckUsers := 24.0 * -1.0

	//iterate over the members_cached array. Maximum limit is 1000.
	for itera := range membersCached {
		userTimeJoin, _ := membersCached[itera].JoinedAt.Parse()
		timevar := userTimeJoin.Sub(time.Now()).Hours()

		if timevar > timeToCheckUsers {
			tempMsg += "This user is too young (less than 24h join age): " + membersCached[itera].User.Username + "\n"
		}
		fmt.Println("Found ", strings.Contains(tempMsg, "\n"), " very young users...")
	}
	//print out the amount of members_cached (max is currently 1000)
	s.ChannelMessageSend(m.ChannelID, tempMsg)
	return

}

// PlanGame Plans a game for a person with a timed reminder
func PlanGame(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	if len(cmd.Arguments) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Insufficient arguments. Provided "+strconv.FormatInt(int64(len(cmd.Arguments)), 10)+" , Expected at least 3")
		return
	}
	GamePlanInsert(&cmd, &s, &m)
	return
}

func PlannedGames(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	if len(cmd.Arguments) > 0 {
		s.ChannelMessageSend(m.ChannelID, "Insufficient arguments. Provided "+strconv.FormatInt(int64(len(cmd.Arguments)), 10)+" , Expected no arguments")
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

	var plannedgames string
	database.DisplayAllGamesPlanned(sqliteDatabase, &plannedgames)

	//send info to channel
	(*s).ChannelMessageSend((*m).ChannelID, plannedgames)
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
		(*s).ChannelMessageSend((*m).ChannelID, "Error parsing time")
		return
	}

	//Put hours into timeHours
	timeHour, err := strconv.Atoi(splitTimeArgument[0])
	if err != nil {
		(*s).ChannelMessageSend((*m).ChannelID, "Error converting hours")
		fmt.Printf("%s", err)
		return
	}
	//put minutes into timeMinute
	timeMinute, err := strconv.Atoi(splitTimeArgument[1])
	if err != nil {
		(*s).ChannelMessageSend((*m).ChannelID, "Error converting minutes")
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
	fmt.Println(t)
	return
}

// GetMemberListFromGuild Gets the member info
func GetMemberListFromGuild(s *discordgo.Session, guildID string) []*discordgo.Member {
	membersList, err := s.GuildMembers(guildID, "0", 1000)
	if err != nil {
		fmt.Println("Error getting the memberlist (probably invalid guildID): " + guildID)
	}

	return membersList

}

// CheckRegularSpamAttack Checks the server for spam attacks
func CheckRegularSpamAttack(s *discordgo.Session) {
	//variable definitons
	var membersCached = GetMemberListFromGuild(s, GuildIDNumber)
	var tempMsg string
	var spamcounter int64
	var checkinterval time.Duration = 90
	var timeToCheckUsers = 10 * -1.0

	for {
		//iterate over the members_cached array. Maximum limit is 1000.
		for itera := range membersCached {
			userTimeJoin, _ := membersCached[itera].JoinedAt.Parse()
			timevar := userTimeJoin.Sub(time.Now()).Minutes()

			if timevar > timeToCheckUsers {
				tempMsg += "RAID PROTECTION ALERT!: User" + membersCached[itera].User.Username + "join age: " + strconv.FormatFloat(timeToCheckUsers, 'f', 0, 64) + "\n"
				spamcounter += 1
			}

		}
		if spamcounter > 4 {
			s.ChannelMessageSend(AdminChannel, "WARN: Possible RAID ATTACK detected!!! (<@&513275201375698954>) ("+strconv.FormatInt(spamcounter, 10)+" users joined in the last "+strconv.FormatFloat(timeToCheckUsers, 'f', 0, 64)+" hours)")
		}
		spamcounter = 0
		time.Sleep(checkinterval * time.Second)
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
			fmt.Println("[ERR] error closing the file with topics")
		}
	}(fileHandle)

	fileScanner := bufio.NewScanner(fileHandle)

	var splitTopic []string

	for fileScanner.Scan() {
		splitTopic = append(splitTopic, fileScanner.Text())
	}

	a := 0
	b := len(splitTopic)

	rand.Seed(time.Now().UnixNano())
	n := a + rand.Intn(b-a+1)

	s.ChannelMessageSend(m.ChannelID, splitTopic[n])
	return
}

func Fox(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "<a:medzernikShake:814055147583438848>")
}

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

	var apiKey = "65bb37a9ac2af9128d6ceaf670043b39"

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

func TimedChannelUnlock(s *discordgo.Session, m *discordgo.MessageCreate) {

}

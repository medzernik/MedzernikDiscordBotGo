// Package responder_functions contains all the logic and basic config commands for the responder commands.
package responder_functions

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/medzernik/SlovakiaDiscordBotGo/command"
	"github.com/medzernik/SlovakiaDiscordBotGo/database"
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
	var muteUser string
	var timeMute uint64
	var roleMuteID = "684159104901709825"
	var priviledgedRole = "513275201375698954"
	membersCached := GetMemberListFromGuild(s, GuildIDNumber)

	if len(cmd.Arguments) < 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "Insufficient arguments provded (help: .mute **@user time**)")
		if err != nil {
			return
		}
		return
	}
	muteUserRaw := cmd.Arguments[0]
	muteUser = command.ParseMentionToString(cmd.Arguments[0])

	if len(cmd.Arguments) > 1 {
		timeMute, err = strconv.ParseUint(cmd.Arguments[1], 10, 64)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "**Error parsing the time** (help: second argument should be number in minutes, max 64-bit sized.) \nExample: .mute user **30**"+"\nYou can also use the command without time specified"+"\n"+err.Error())
			return
		}
	}

	for itera := range membersCached {
		if membersCached[itera].User.ID == muteUser {
			for i := range membersCached[itera].Roles {
				if membersCached[itera].Roles[i] == roleMuteID {
					s.ChannelMessageSend(m.ChannelID, "**Duplicate request:** User "+muteUser+" already has a role Muted "+roleMuteID+" added")
					return
				} else if muteUser == "837982234597916672" {
					s.ChannelMessageSend(m.ChannelID, "***You fucking donkey. WHO DO YOU THINK YOU ARE MUTING YOU FUCK*** <:pocem:683715405407059968> <:pocem:683715405407059968> <:pocem:683715405407059968> <:pocem:683715405407059968>")
					return
				} else if membersCached[itera].Roles[i] == priviledgedRole && membersCached[itera].User.ID == muteUser {
					s.ChannelMessageSend(m.ChannelID, "**Refusing to mute user** "+muteUserRaw+" **who has a privileged role** "+priviledgedRole+"\n The role is set as priviledged and can't be assigned role id: "+roleMuteID)
					return

				}

			}

		}
	}
	for itera := range membersCached {
		for i := range membersCached[itera].Roles {
			if membersCached[itera].User.ID == m.Author.ID && membersCached[itera].Roles[i] == priviledgedRole {
				var err error
				err = s.GuildMemberMute(GuildIDNumber, muteUser, true)
				if err != nil {
				}
				err = s.GuildMemberRoleAdd(GuildIDNumber, muteUser, roleMuteID)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "Error adding muted role to user "+muteUserRaw+"\n"+err.Error())
					return
				}
				s.ChannelMessageSend(m.ChannelID, "User muted successfully.")
				s.ChannelMessageSend(LogChannel, "User "+muteUserRaw+" muted successfully - assigned role "+roleMuteID)
				return
			} else {
				s.ChannelMessageSend(m.ChannelID, "Insufficient permissions")
				fmt.Println(timeMute)
				return
			}

		}
	}
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

}

// PlanGame Plans a game for a person with a timed reminder
func PlanGame(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	if len(cmd.Arguments) < 3 {
		s.ChannelMessageSend(m.ChannelID, "Insufficient arguments. Provided "+strconv.FormatInt(int64(len(cmd.Arguments)), 10)+" , Expected at least 3")
		return
	}
	GamePlanInsert(&cmd, &s, &m)
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
				println("Checking " + membersCached[itera].User.Username + " for possible raid attack")
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

func Trivia(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	err := command.VerifyArguments(&cmd)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	fileHandle, err := os.Open("trivia_questions.txt")
	if err != nil {
		fmt.Println("error reading the file: ", err)
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	defer func(fileHandle *os.File) {
		err := fileHandle.Close()
		if err != nil {
			fmt.Println("error closing the file with trivia")
		}
	}(fileHandle)

	fileScanner := bufio.NewScanner(fileHandle)

	var splitTrivia []string

	for fileScanner.Scan() {
		splitTrivia = append(splitTrivia, fileScanner.Text())
	}

	a := 0
	b := len(splitTrivia)

	rand.Seed(time.Now().UnixNano())
	n := a + rand.Intn(b-a+1)

	s.ChannelMessageSend(m.ChannelID, splitTrivia[n])

}

func Fox(s *discordgo.Session, cmd command.Command, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, "<a:medzernikShake:814055147583438848>")
}

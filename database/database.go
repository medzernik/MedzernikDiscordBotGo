package database

import (
	"database/sql"
	"fmt"
	_ "fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	_ "io/ioutil"
	"log"
	_ "modernc.org/sqlite"
	"os"
	_ "path/filepath"
	"strconv"
	"time"
)

// Databaserun will delete the old database and then create a new one, get all the file handlers and basic info
func Databaserun() {
	var test string

	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db") // Open the created SQLite File
	defer sqliteDatabase.Close()                                     // Defer Closing the database
	createTable(sqliteDatabase)                                      // Create Database Tables

	// DISPLAY INSERTED RECORDS
	DisplayGamePlanned(sqliteDatabase, &test)
}

// createTable creates a game planning table
func createTable(db *sql.DB) {
	createGamePlanningDB := `CREATE TABLE gameplanning (
		"idGames" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"time" INTEGER,
		"gamename" TEXT,
		"mentions" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create game table...")
	statement, err := db.Prepare(createGamePlanningDB) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("game table created")
}

// InsertGame inserts the requested data into the database
func InsertGame(db *sql.DB, time int64, gamename string, mentions string) {
	log.Println("Inserting game record ...")
	insertGamePlanning := `INSERT INTO gameplanning(time, gamename, mentions) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertGamePlanning) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(time, gamename, mentions)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// DisplayGamePlanned will immidiately display the currently planned game into the channel it was planned at as a confirmation for the user. Returns a string with all the neccessary data
func DisplayGamePlanned(db *sql.DB, output *string) string {
	row, err := db.Query("SELECT * FROM gameplanning ORDER BY gamename")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var timestamp int64
		var gamename string
		var mentions string
		row.Scan(&id, &timestamp, &gamename, &mentions)
		log.Println("Game is planned for: ", timestamp, " ", gamename, " ", mentions)
		*output = "ID: " + strconv.FormatInt(int64(id), 10) + ", cas: " + time.Unix(timestamp, 0).Format(time.RFC822) + ", hra: " + gamename + ", s ludmi " + mentions + "\n"
	}
	return *output
}

// DisplayAllGamesPlanned displays all planned games in the database and outputs the result into the channel.
func DisplayAllGamesPlanned(db *sql.DB, output *string) string {
	row, err := db.Query("SELECT * FROM gameplanning ORDER BY time")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var timestamp int64
		var gamename string
		var mentions string
		row.Scan(&id, &timestamp, &gamename, &mentions)
		log.Println("Game is planned for: ", timestamp, " ", gamename, " ", mentions)
		*output += "ID: " + strconv.FormatInt(int64(id), 10) + ", cas: " + time.Unix(timestamp, 0).Format(time.RFC822) + ", hra: " + gamename + ", s ludmi " + mentions + "\n"
	}
	return *output
}

// CheckPlannedGames runs concurrently with the go command at bot startup.
func CheckPlannedGames(s **discordgo.Session) {
	var checkInterval time.Duration = 59
	//This is here for the function to wait until the database is created (since it's async). I should *really* make this a proper way, not a fixed wait time...)
	var initInterval time.Duration = 2
	//Channel into which to output the information
	var gameReminderChannelID string = "837987736416813076"

	fmt.Println("Initializing CheckPlannedGames module")
	time.Sleep(initInterval * time.Second)
	fmt.Println("CheckPlannedGames module initialized successfully...")

	//Loop that continuously runs... With a timer to wait for 59 seconds
	for {
		sqliteDatabase, _ := sql.Open("sqlite3", "./sqlite-database.db")
		plannedGame, err := sqliteDatabase.Query("SELECT * FROM gameplanning ORDER BY time")
		if err != nil {
			log.Fatal(err)
		}

		for plannedGame.Next() {
			var id int
			var timestamp int64
			var gamename string
			var mentions string
			plannedGame.Scan(&id, &timestamp, &gamename, &mentions)

			var timestampInt = time.Unix(timestamp, 0)

			if time.Now().Date() == timestampInt.Date() && time.Now().Hour() == timestampInt.Hour() && time.Now().Minute() == timestampInt.Minute() {
				(*s).ChannelMessageSend(gameReminderChannelID, "**CAS SA HRAT** "+"cas: "+time.Unix(timestamp, 0).Format(time.RFC822)+", hra: "+gamename+", s ludmi "+mentions+"\n")
			}

		}
		//Release the database
		sqliteDatabase.Close()
		//Wait until checking again for 59 seconds (checkInterval)
		time.Sleep(checkInterval * time.Second)
	}
}

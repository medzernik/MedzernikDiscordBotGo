package database

import (
	"database/sql"
	_ "fmt"
	_ "github.com/mattn/go-sqlite3"
	_ "io/ioutil"
	"log"
	_ "modernc.org/sqlite"
	"os"
	_ "path/filepath"
)

func Databaserun() {
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

	// INSERT RECORDS
	insertGame(sqliteDatabase, "12:45", "Terraria", "medzernik")

	// DISPLAY INSERTED RECORDS
	displayGamePlanning(sqliteDatabase)
}
func createTable(db *sql.DB) {
	createGamePlanningDB := `CREATE TABLE gameplanning (
		"idGames" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"time" TEXT,
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

// We are passing db reference connection from main to our method with other parameters
func insertGame(db *sql.DB, time string, gamename string, mentions string) {
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

func displayGamePlanning(db *sql.DB) {
	row, err := db.Query("SELECT * FROM gameplanning ORDER BY gamename")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var time string
		var gamename string
		var mentions string
		row.Scan(&id, &time, &gamename, &mentions)
		log.Println("Game is planned for: ", time, " ", gamename, " ", mentions)
	}
}

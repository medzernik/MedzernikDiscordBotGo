// Package database This serves to save a database of the members so we don't have to constantly get it from somewhere. Will be replaced by a redis cache.
package database

/*
import (
	"database/sql"
	"fmt"
	"log"
)

func createBirthdayDatabase(db *sql.DB) {
	createMembersDB := `CREATE TABLE birthday (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"time" INTEGER,
        "metnion" TEXT
	  );` // SQL Statement for Create Table

	log.Println("Create game table...")
	statement, err := db.Prepare(createMembersDB) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err3 := statement.Exec()
	if err3 != nil {
		fmt.Println("error executing SQL statements")
		return
	} // Execute SQL Statements
	log.Println("game table created")

}

*/

package logging

import (
	"log"
	"os"
)

var (
	LoggerWarning *log.Logger
	LoggerInfo    *log.Logger
	LoggerError   *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	LoggerInfo = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	LoggerWarning = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	LoggerError = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	LoggerInfo.Println("Starting bot...")
}

package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func initLogger() {
	now := time.Now()
	fileName := fmt.Sprintf("%v.log", now.Format("2006-01-02-15"))

	file, err := os.OpenFile("runtime/logs/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	fileWarning, err := os.OpenFile("runtime/warning/warning_"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open log error file:", err)
	}

	fileErr, err := os.OpenFile("runtime/errors/error_"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("Failed to open log error file:", err)
	}

	Info = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(fileWarning, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(fileErr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func CheckAndRotateLog() {
	initLogger()
}

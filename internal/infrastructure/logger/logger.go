package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func InitLogger() {
	currentDate := time.Now().Format("2006-01-02")

	// pastikan folder logs ada
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		log.Fatal("Failed to create log directory:", err)
	}

	// nama file
	logInfoFileName := fmt.Sprintf("logs/info-%s.log", currentDate)
	logErrorFileName := fmt.Sprintf("logs/error-%s.log", currentDate)

	// buka file info
	infoFile, err := os.OpenFile(logInfoFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to open info log file:", err)
	}

	// buka file error
	errorFile, err := os.OpenFile(logErrorFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to open error log file:", err)
	}

	// assign logger
	InfoLogger = log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

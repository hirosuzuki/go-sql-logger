package logger

import (
	"fmt"
	"log"
	"os"
)

var sqlLogFileName string
var sqlLogFile *os.File

func init() {
	sqlLogFileName = os.Getenv("SQL_LOGFILE")
	if sqlLogFileName == "" {
		sqlLogFileName = "/tmp/sql.log"
	}

	log.Printf("Go SQL Logger: Log File -> %s\n", sqlLogFileName)

	var err error
	if sqlLogFile, err = os.OpenFile(sqlLogFileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644); err != nil {
		log.Printf("Go SQL Logger: Open File Error: %s\n", err.Error())
		return
	}
}

func Log(startTime int64, timeDelta int64, tag string, query string) {
	if sqlLogFile == nil {
		return
	}
	fmt.Fprintf(sqlLogFile, "%d\t%d\t%s\t%s\n", startTime, timeDelta, tag, query)
}

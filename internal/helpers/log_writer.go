package helpers

import (
	"log"
	"os"
)

var NzLogInfo *log.Logger
var NzLogError *log.Logger

func InitNzLog() {
	// Make sure log_full_path exists
	if err := os.MkdirAll(Conf.LogDir, 0770); err != nil {
		log.Fatalln("Failed to create log full path recursively: ", err)
	}

	fl, err := os.OpenFile(Conf.LogDir+"cron_backup_log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0770)
	if err != nil {
		log.Fatalln("Failed to open|create log file: ", err)
	}

	NzLogInfo = log.New(fl, "[INFO] ", log.Ldate|log.Ltime)
	NzLogError = log.New(fl, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}

package utils

import (
	"log"
	"os"
	"os/user"
)

// testCreateDir test to create log dir and backup dir
func testCreateDir() bool {
	isPass := true

	usr, _ := user.Current()
	log.Println("[INFO] Try creating log dir and backup dir with user:", usr.Username)

	if err := os.MkdirAll(testConf.LogDir, 0770); err != nil {
		log.Println("[ERROR] Failed to create log dir:", err)
		isPass = false
	}
	if err := os.MkdirAll(testConf.Backup.RootDir, 0770); err != nil {
		log.Println("[ERROR] Failed to create log dir:", err)
		isPass = false
	}
	if err := os.MkdirAll(testConf.BackupAppDir, 0770); err != nil {
		log.Println("[ERROR] Failed to create backup app dir:", err)
		isPass = false
	}
	if err := os.MkdirAll(testConf.BackupDBDir, 0770); err != nil {
		log.Println("[ERROR] Failed to create backup db dir:", err)
		isPass = false
	}

	return isPass
}

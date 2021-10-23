package utils

import (
	"log"

	"github.com/mdanialr/go-cron-backup/internal/models"
)

var testConf *models.Config
var fileToDelete struct {
	APPname string
	DBname  string
}

// RunTest main function to run all checking and testing then throw all errors if any
func RunTest(isDel bool, isExDB bool) {
	if isPass := testCheckConfig(); isPass {
		log.Println("[INFO] Success checking config file")
	}
	if isPass := testCreateDir(); isPass {
		log.Println("[INFO] Success creating log and backup folder")
	}
	if isPass := testBackup(isExDB); isPass {
		log.Println("[INFO] Success creating backup for database and app")
	}
	if isPass := testDelete(isDel); isPass {
		log.Println("[INFO] Success deleting dir that created by this test")
	}
	log.Println("[INFO] Successfully testing all config file and functionality!")
	log.Println("[INFO] Next is you can create cronjob to run this app many times as needed")
}

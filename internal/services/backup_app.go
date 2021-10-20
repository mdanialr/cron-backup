package services

import (
	"log"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
	"github.com/mdanialr/go-cron-backup/internal/models"
)

// backupAPP do the backup using goroutine
func backupAPP(wg *sync.WaitGroup) {
	for _, v := range helpers.Conf.Backup.APP.Apps {
		backupDir := helpers.Conf.BackupAppDir + v.App.DirName
		if err := makeSureDirExists(backupDir); err != nil {
			log.Fatalf("Failed to create dir for backup app in %v: %v\n", v.App.AppDir, err)
		}
		commands := parseBackupAPPCommand(v.App)

		// delete old backup
		deleteOldBackupAPP(wg, backupDir)

		// goroutine to separate zip proccess from main thread
		wg.Add(1)
		go func(wg *sync.WaitGroup, app models.App) {
			helpers.NzLogInfo.Println("[START] zipping in", "'"+app.AppDir+"'")
			out, err := exec.Command("sh", "-c", commands).CombinedOutput()
			if err != nil {
				helpers.NzLogError.Println(string(out))
			}
			helpers.NzLogInfo.Println("[DONE] zipping", "'"+app.DirName+"'")

			wg.Done()
		}(wg, v.App)
	}
	wg.Done()
}

// parseBackupAPPCommand combine all commands and args
func parseBackupAPPCommand(app models.App) string {
	fmtTime := time.Now().Format("2006-Jan-02_Monday_15:04:05")
	fName := "/" + fmtTime + ".zip"
	zipName := helpers.Conf.BackupAppDir + app.DirName + fName
	cmdSeries := []string{
		"cd " + app.AppDir,
		"zip -r -q " + zipName + " *",
	}
	return strings.Join(cmdSeries, ";")
}

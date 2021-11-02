package services

import (
	"log"
	"sync"
	"time"

	"github.com/mdanialr/go-cron-backup/internal/arch"
	"github.com/mdanialr/go-cron-backup/internal/helpers"
	"github.com/mdanialr/go-cron-backup/internal/models"
)

// backupAPP deleting old backup then zipping dir.
func backupAPP(wg *sync.WaitGroup) {
	defer wg.Done()

	for _, v := range helpers.Conf.Backup.APP.Apps {
		backupDir := helpers.Conf.BackupAppDir + v.App.DirName
		if err := makeSureDirExists(backupDir); err != nil {
			log.Fatalf("Failed to create dir for backup app in %v: %v\n", v.App.AppDir, err)
		}

		// delete old backup
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			if err := deleteOldBackup(backupDir, helpers.Conf.Backup.APP.Retain); err != nil {
				helpers.NzLogError.Println(err)
			}
		}(wg)

		// goroutine to separate zip proccess from main thread
		wg.Add(1)
		go func(wg *sync.WaitGroup, app models.App) {
			defer wg.Done()

			helpers.NzLogInfo.Println("[START] zipping in", "'"+app.AppDir+"'")

			fmtTime := time.Now().Format("2006-Jan-02_Monday_15:04:05")
			fName := "/" + fmtTime + ".zip"
			zipName := helpers.Conf.BackupAppDir + app.DirName + fName

			if err := arch.ZipDir(app.AppDir, zipName); err != nil {
				helpers.NzLogError.Printf("Failed zipping in %v: %v", app.DirName, err)
			}

			helpers.NzLogInfo.Println("[DONE] zipping", "'"+app.DirName+"'")
		}(wg, v.App)
	}
}

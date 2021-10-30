package services

import (
	"log"
	"sync"

	"github.com/mdanialr/go-cron-backup/internal/arch"
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

		// delete old backup
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			if err := deleteOldBackup(backupDir, helpers.Conf.Backup.APP.Retain); err != nil {
				helpers.NzLogError.Println(err)
			}
			wg.Done()
		}(wg)

		// goroutine to separate zip proccess from main thread
		wg.Add(1)
		go func(wg *sync.WaitGroup, app models.App) {
			helpers.NzLogInfo.Println("[START] zipping in", "'"+app.AppDir+"'")
			if err := arch.BashAPPZip(app); err != nil {
				helpers.NzLogError.Println(err)
			}
			helpers.NzLogInfo.Println("[DONE] zipping", "'"+app.DirName+"'")

			wg.Done()
		}(wg, v.App)
	}
	wg.Done()
}

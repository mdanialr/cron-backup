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

	// initialize number of jobs and the job channel
	numJobs := len(helpers.Conf.Backup.APP.Apps)
	jobChan := make(chan models.App, numJobs)
	doneChan := make(chan int, numJobs)

	// start workers.
	for w := 1; w <= helpers.Conf.Backup.APP.MaxWorker; w++ {
		go appWorker(jobChan, doneChan)
	}

	// send jobs before closing the sending channel.
	for _, v := range helpers.Conf.Backup.APP.Apps {
		jobChan <- v.App
	}
	close(jobChan)

	for range helpers.Conf.Backup.APP.Apps {
		// block until all jobs is done.
		<-doneChan
	}
}

// dbWorker worker function to do the job which is
// deleting old backup, and zipping app dir or folder.
func appWorker(jobChan <-chan models.App, doneChan chan<- int) {
	// listen to job channel.
	for app := range jobChan {
		// just send whatever number to channel.
		doneChan <- 1

		// make sure target backup dir is exist by creating it.
		backupDir := helpers.Conf.BackupAppDir + app.DirName
		if err := makeSureDirExists(backupDir); err != nil {
			log.Fatalf("Failed to create dir for backup app in %v: %v\n", app.AppDir, err)
		}

		// delete old backup, according to maximum retain days
		// in config.
		if err := deleteOldBackup(backupDir, helpers.Conf.Backup.APP.Retain); err != nil {
			helpers.NzLogError.Println(err)
		}

		// zipping app directory or folders.
		helpers.NzLogInfo.Println("[START] zipping in", "'"+app.AppDir+"'")

		fmtTime := time.Now().Format("2006-Jan-02_Monday_15:04:05")
		fName := "/" + fmtTime + ".zip"
		zipName := helpers.Conf.BackupAppDir + app.DirName + fName

		if err := arch.ZipDir(app.AppDir, zipName); err != nil {
			helpers.NzLogError.Printf("Failed zipping in %v: %v", app.DirName, err)
		}

		helpers.NzLogInfo.Println("[DONE] zipping", "'"+app.DirName+"'")
	}
}

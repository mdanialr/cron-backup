package services

import (
	"log"
	"os/exec"
	"sync"
	"time"

	"github.com/mdanialr/go-cron-backup/internal/arch"
	"github.com/mdanialr/go-cron-backup/internal/helpers"
	"github.com/mdanialr/go-cron-backup/internal/models"
)

// backupDB deleting old backup, dump database then zip them.
func backupDB(wg *sync.WaitGroup) {
	defer wg.Done()

	// initialize number of jobs and the job channel
	numJobs := len(helpers.Conf.Backup.DB.Databases)
	jobChan := make(chan models.Database, numJobs)
	doneChan := make(chan int, numJobs)

	// start workers.
	for w := 1; w <= 2; w++ {
		go dbWorker(jobChan, doneChan)
	}

	// send jobs.
	for _, v := range helpers.Conf.Backup.DB.Databases {
		jobChan <- v.Database
	}
	close(jobChan)

	for range helpers.Conf.Backup.DB.Databases {
		// block until all jobs is done.
		<-doneChan
	}
}

// dbWorker worker function to do the job which is
// deleting old backup, dumping database, and zip them.
func dbWorker(jobChan <-chan models.Database, doneChan chan<- int) {
	// listen to job channel.
	for db := range jobChan {
		// just send whatever number to channel.
		doneChan <- 1

		// make sure target backup dir is exist by creating it.
		backupDir := helpers.Conf.BackupDBDir + db.DirName
		if err := makeSureDirExists(backupDir); err != nil {
			log.Fatalf("Failed to create dir for backup app in %v: %v\n", db.DirName, err)
		}

		// setup dump command and the dumped output file name.
		var dumpCmd, outName string
		if db.T.MariaDB {
			dumpCmd, outName = parseDumpingMariaDBCommand(db)
		}
		if db.T.PGsql {
			dumpCmd, outName = parseDumpingPGCommand(db)
		}

		// delete old backup, according to maximum retain days
		// in config.
		if err := deleteOldBackup(backupDir, helpers.Conf.Backup.DB.Retain); err != nil {
			helpers.NzLogError.Println(err)
		}

		// dumping database
		helpers.NzLogInfo.Println("[START] dumping database", "'"+db.Name+"'")
		if _, err := exec.Command("sh", "-c", dumpCmd).CombinedOutput(); err != nil {
			helpers.NzLogError.Println(err)
		}
		helpers.NzLogInfo.Println("[DONE] dumping", "'"+db.Name+"'")

		// zipping dumped database
		helpers.NzLogInfo.Println("[START] zipping dumped database", "'"+db.Name+"'")

		fmtTime := time.Now().Format("2006-Jan-02_Monday_15:04:05")
		fName := "/" + fmtTime + ".zip"
		zipName := helpers.Conf.BackupDBDir + db.DirName + fName

		if err := arch.Zip("/tmp/"+outName, zipName); err != nil {
			helpers.NzLogError.Printf("Failed zipping %v: %v", outName, err)
		}

		helpers.NzLogInfo.Println("[DONE] zipping", "'"+db.Name+"'")
	}
}

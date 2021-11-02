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

// backupDB do the backup using goroutine
func backupDB(wg *sync.WaitGroup) {
	defer wg.Done()

	for _, v := range helpers.Conf.Backup.DB.Databases {
		backupDir := helpers.Conf.BackupDBDir + v.Database.DirName
		if err := makeSureDirExists(backupDir); err != nil {
			log.Fatalf("Failed to create dir for backup app in %v: %v\n", v.Database.DirName, err)
		}

		var dumpCmd, outName string
		if v.Database.T.MariaDB {
			dumpCmd, outName = parseDumpingMariaDBCommand(v.Database)
		}
		if v.Database.T.PGsql {
			dumpCmd, outName = parseDumpingPGCommand(v.Database)
		}

		// delete old backup
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			if err := deleteOldBackup(backupDir, helpers.Conf.Backup.DB.Retain); err != nil {
				helpers.NzLogError.Println(err)
			}
		}(wg)

		// goroutine to separate zip proccess from main thread
		wg.Add(1)
		go func(wg *sync.WaitGroup, db models.Database) {
			defer wg.Done()
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
		}(wg, v.Database)
	}
}

package services

import (
	"log"
	"os/exec"
	"strings"
	"sync"

	"github.com/mdanialr/go-cron-backup/internal/arch"
	"github.com/mdanialr/go-cron-backup/internal/helpers"
	"github.com/mdanialr/go-cron-backup/internal/models"
)

// backupDB do the backup using goroutine
func backupDB(wg *sync.WaitGroup) {
	for _, v := range helpers.Conf.Backup.DB.Databases {
		backupDir := helpers.Conf.BackupDBDir + v.Database.DirName
		if err := makeSureDirExists(backupDir); err != nil {
			log.Fatalf("Failed to create dir for backup app in %v: %v\n", v.Database.DirName, err)
		}
		dumpCmd, outName := parseDumpingDBCmd(v.Database)

		// delete old backup
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			if err := deleteOldBackup(backupDir, helpers.Conf.Backup.DB.Retain); err != nil {
				helpers.NzLogError.Println(err)
			}
			wg.Done()
		}(wg)

		// goroutine to separate zip proccess from main thread
		wg.Add(1)
		go func(wg *sync.WaitGroup, db models.Database) {
			// dumping database
			helpers.NzLogInfo.Println("[START] dumping database", "'"+db.Name+"'")
			if _, err := exec.Command("sh", "-c", dumpCmd).CombinedOutput(); err != nil {
				helpers.NzLogError.Println(err)
			}
			helpers.NzLogInfo.Println("[DONE] dumping", "'"+db.Name+"'")

			// zipping dumped database
			helpers.NzLogInfo.Println("[START] zipping dumped database", "'"+db.Name+"'")
			if err := arch.BashDBZip(db, outName); err != nil {
				helpers.NzLogError.Println(err)
			}
			helpers.NzLogInfo.Println("[DONE] zipping", "'"+db.Name+"'")

			wg.Done()
		}(wg, v.Database)
	}

	wg.Done()
}

// parseDumpingDBCmd combine all commands for dumping database
func parseDumpingDBCmd(db models.Database) (string, string) {
	cmd := "sudo -u postgres pg_dump"
	args := "--clean --no-owner"
	outName := "dump_" + db.Name
	cmdSeries := []string{
		cmd,
		db.Name,
		args,
		">",
		outName,
	}
	dumpCmd := strings.Join(cmdSeries, " ")
	return strings.Join([]string{"cd /tmp", dumpCmd}, ";"), outName
}

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

// backupDB do the backup using goroutine
func backupDB(wg *sync.WaitGroup) {
	for _, v := range helpers.Conf.Backup.DB.Databases {
		backupDir := helpers.Conf.BackupDBDir + v.Database.DirName
		if err := makeSureDirExists(backupDir); err != nil {
			log.Fatalf("Failed to create dir for backup app in %v: %v\n", v.Database.DirName, err)
		}
		dumpCmd, outName := parseDumpingCommand(v.Database)
		zipCmd := parseZippingCommand(v.Database, outName)

		// delete old backup
		deleteOldBackupDB(wg, backupDir)

		// goroutine to separate zip proccess from main thread
		wg.Add(1)
		go func(wg *sync.WaitGroup, db models.Database) {
			// dumping database
			helpers.NzLogInfo.Println("[START] dumping database", "'"+db.Name+"'")
			out, err := exec.Command("sh", "-c", dumpCmd).CombinedOutput()
			if err != nil {
				helpers.NzLogError.Println(string(out))
			}
			helpers.NzLogInfo.Println("[DONE] dumping", "'"+db.Name+"'")

			// zipping dumped database
			helpers.NzLogInfo.Println("[START] zipping dumped database", "'"+db.Name+"'")
			out, err = exec.Command("sh", "-c", zipCmd).CombinedOutput()
			if err != nil {
				helpers.NzLogError.Println(string(out))
			}
			helpers.NzLogInfo.Println("[DONE] zipping", "'"+db.Name+"'")

			wg.Done()
		}(wg, v.Database)
	}
	wg.Done()
}

// parseDumpingCommand combine all commands for dumping database
func parseDumpingCommand(db models.Database) (string, string) {
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

// parseZippingCommand combine all commands for zipping dumped database
func parseZippingCommand(db models.Database, outName string) string {
	fmtTime := time.Now().Format("2006-Jan-02_Monday_15:04:05")
	fName := "/" + fmtTime + ".zip"
	zipName := helpers.Conf.BackupDBDir + db.DirName + fName
	cmdSeries := []string{
		"cd /tmp",
		"zip -q " + zipName + " " + outName,
	}
	return strings.Join(cmdSeries, ";")
}

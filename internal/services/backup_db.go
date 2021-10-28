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
		dumpCmd, outName := parseDumpingDBCmd(v.Database)
		zipCmd := parseZippingDBCmd(v.Database, outName)

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
			out, err := exec.Command("sh", "-c", dumpCmd).CombinedOutput()
			if err != nil {
				helpers.NzLogError.Println(string(out))
				helpers.NzLogError.Println(err)
			}
			helpers.NzLogInfo.Println("[DONE] dumping", "'"+db.Name+"'")

			// zipping dumped database
			helpers.NzLogInfo.Println("[START] zipping dumped database", "'"+db.Name+"'")
			out, err = exec.Command("sh", "-c", zipCmd).CombinedOutput()
			if err != nil {
				helpers.NzLogError.Println(string(out))
				helpers.NzLogError.Println(err)
			}
			helpers.NzLogInfo.Println("[DONE] zipping", "'"+db.Name+"'")

			wg.Done()
		}(wg, v.Database)
	}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		if err := deleteDumpedFile(); err != nil {
			helpers.NzLogError.Println(err)
		}
		wg.Done()
	}(wg)

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

// parseZippingDBCmd combine all commands for zipping dumped database
func parseZippingDBCmd(db models.Database, outName string) string {
	fmtTime := time.Now().Format("2006-Jan-02_Monday_15:04:05")
	fName := "/" + fmtTime + ".zip"
	zipName := helpers.Conf.BackupDBDir + db.DirName + fName
	cmdSeries := []string{
		"cd /tmp",
		"zip -q " + zipName + " " + outName,
	}
	return strings.Join(cmdSeries, ";")
}

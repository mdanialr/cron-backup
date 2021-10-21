package utils

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mdanialr/go-cron-backup/internal/models"
)

// testBackup try to run backup and handle the goroutine
func testBackup() bool {
	isPass := true

	cAPP := make(chan bool)
	cDB := make(chan bool)

	go testBackupAPP(cAPP)
	go testBackupDB(cDB)

	if !<-cAPP || !<-cDB {
		isPass = false
	}

	return isPass
}

// testBackupAPP try to run backup on first database in config file
func testBackupAPP(c chan bool) {
	isPass := true

	tAPP := testConf.Backup.APP.Apps[0].App
	backupDir := testConf.BackupAppDir + tAPP.DirName
	if err := os.MkdirAll(backupDir, 0770); err != nil {
		log.Fatalf("Failed to create dir for backup app in %v: %v\n", tAPP.AppDir, err)
	}
	commands := parseBackupAPPCommand(tAPP)

	log.Println("[START] zipping in", "'"+tAPP.AppDir+"'")
	out, err := exec.Command("sh", "-c", commands).CombinedOutput()
	if err != nil {
		log.Println(string(out))
	}
	log.Println("[DONE] zipping", "'"+tAPP.DirName+"'")

	c <- isPass
}

// parseBackupAPPCommand combine all commands and args
func parseBackupAPPCommand(app models.App) string {
	fmtTime := time.Now().Format("2006-Jan-02_Monday_15:04:05")
	fName := "/" + fmtTime + ".zip"
	zipName := testConf.BackupAppDir + app.DirName + fName
	cmdSeries := []string{
		"cd " + app.AppDir,
		"zip -r -q " + zipName + " *",
	}

	return strings.Join(cmdSeries, ";")
}

// testBackupDB try to run backup on first database in config file
func testBackupDB(c chan bool) {
	isPass := true

	tDB := testConf.Backup.DB.Databases[0].Database
	backupDir := testConf.BackupDBDir + tDB.DirName
	if err := os.MkdirAll(backupDir, 0770); err != nil {
		log.Fatalf("Failed to create dir for backup app in %v: %v\n", tDB.DirName, err)
	}
	dumpCmd, outName := parseDumpingCommand(tDB)
	zipCmd := parseZippingCommand(tDB, outName)

	// dumping database
	log.Println("[START] dumping database", "'"+tDB.Name+"'")
	out, err := exec.Command("sh", "-c", dumpCmd).CombinedOutput()
	if err != nil {
		log.Println(string(out))
		isPass = false
	}
	log.Println("[DONE] dumping", "'"+tDB.Name+"'")

	// zipping dumped database
	log.Println("[START] zipping dumped database", "'"+tDB.Name+"'")
	out, err = exec.Command("sh", "-c", zipCmd).CombinedOutput()
	if err != nil {
		log.Println(string(out))
		isPass = false
	}
	log.Println("[DONE] zipping", "'"+tDB.Name+"'")

	c <- isPass
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
	zipName := testConf.BackupDBDir + db.DirName + fName
	cmdSeries := []string{
		"cd /tmp",
		"zip -q " + zipName + " " + outName,
	}

	return strings.Join(cmdSeries, ";")
}

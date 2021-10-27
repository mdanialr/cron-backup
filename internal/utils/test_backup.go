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
func testBackup(isExDB bool, isExAPP bool) bool {
	isPass := true

	cAPP := make(chan bool)
	cDB := make(chan bool)
	defer close(cAPP)
	defer close(cDB)

	if isExAPP {
		log.Println("[INFO] Excluding app from this test")
	}
	if isExDB {
		log.Println("[INFO] Excluding database from this test")
	}
	if !isExAPP {
		go testBackupAPP(cAPP)
		if !<-cAPP {
			isPass = false
		}
	}
	if !isExDB {
		go testBackupDB(cDB)
		if !<-cDB {
			isPass = false
		}
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
	fileToDelete.APPname = zipName

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

	var dumpCmd, outName string
	if tDB.T.MariaDB {
		dumpCmd, outName = parseDumpingMariaDBCommand(tDB)
	}
	if tDB.T.PGsql {
		dumpCmd, outName = parseDumpingPGCommand(tDB)
	}
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

	// delete dumped database from /tmp
	if err := testDeleteDumpedFile(); err != nil {
		log.Println(string(out))
		isPass = false
	}

	c <- isPass
}

// parseDumpingMariaDBCommand combine all commands for dumping database
func parseDumpingMariaDBCommand(db models.Database) (string, string) {
	cmd := "mariadb-dump " + db.Name
	usr := "-u " + db.Usr
	pwd := "-p" + db.Pwd
	outName := "dump_" + db.Name
	cmdSeries := []string{
		cmd,
		usr,
		pwd,
		">",
		outName,
	}
	dumpCmd := strings.Join(cmdSeries, " ")
	return strings.Join([]string{"cd /tmp", dumpCmd}, ";"), outName
}

// parseDumpingPGCommand combine all commands for dumping database
func parseDumpingPGCommand(db models.Database) (string, string) {
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
	fileToDelete.DBname = zipName

	return strings.Join(cmdSeries, ";")
}

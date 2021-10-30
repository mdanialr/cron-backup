package utils

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mdanialr/go-cron-backup/internal/arch"
	"github.com/mdanialr/go-cron-backup/internal/helpers"
	"github.com/mdanialr/go-cron-backup/internal/models"
)

// testBackup try to run backup and handle the goroutine
func testBackup() bool {
	isPass := true

	cAPP := make(chan bool)
	cDB := make(chan bool)

	if helpers.TCond.IsNoAPP {
		log.Println("[INFO] Excluding app from this test")
	}
	if helpers.TCond.IsNoDB {
		log.Println("[INFO] Excluding database from this test")
	}
	if !helpers.TCond.IsNoAPP {
		go testBackupAPP(cAPP)
		if !<-cAPP {
			isPass = false
		}
	}
	if !helpers.TCond.IsNoDB {
		go testBackupDB(cDB)
		if !<-cDB {
			isPass = false
		}
	}
	if !isPass {
		os.Exit(1)
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

	log.Println("[START] zipping in", "'"+tAPP.AppDir+"'")

	fmtTime := time.Now().Format("2006-Jan-02_Monday_15:04:05")
	fName := "/" + fmtTime + ".zip"
	zipName := testConf.BackupAppDir + tAPP.DirName + fName

	if err := arch.ZipDir(tAPP.AppDir, zipName); err != nil {
		log.Fatalf("Failed zipping in %v: %v", tAPP.DirName, err)
	}

	log.Println("[DONE] zipping", "'"+tAPP.DirName+"'")

	c <- isPass
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

	fmtTime := time.Now().Format("2006-Jan-02_Monday_15:04:05")
	fName := "/" + fmtTime + ".zip"
	zipName := testConf.BackupDBDir + tDB.DirName + fName

	if err := arch.Zip("/tmp/"+outName, zipName); err != nil {
		log.Fatalf("Failed zipping %v: %v", outName, err)
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

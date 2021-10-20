package services

import (
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
)

// deleteOldBackupAPP delete old backup using goroutine
func deleteOldBackupAPP(wg *sync.WaitGroup, dir string) {
	commands := parseDelOldBackupAPPCommand(dir)

	// goroutine to separate delete proccess from main thread
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		helpers.NzLogInfo.Println("[START] deleting old backup in", "'"+dir+"'")
		out, err := exec.Command("sh", "-c", commands).CombinedOutput()
		if err != nil {
			helpers.NzLogError.Println(string(out))
		}
		helpers.NzLogInfo.Println("[DONE] deleting old backup in", "'"+dir+"'")

		wg.Done()
	}(wg)
}

// deleteOldBackupDB delete old backup using goroutine
func deleteOldBackupDB(wg *sync.WaitGroup, dir string) {
	commands := parseDelOldBackupDBCommand(dir)

	// goroutine to separate delete proccess from main thread
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		helpers.NzLogInfo.Println("[START] deleting old backup in", "'"+dir+"'")
		out, err := exec.Command("sh", "-c", commands).CombinedOutput()
		if err != nil {
			helpers.NzLogError.Println(string(out))
		}
		helpers.NzLogInfo.Println("[DONE] deleting old backup in", "'"+dir+"'")

		wg.Done()
	}(wg)
}

// parseDelOldBackupAPPCommand combine all args
func parseDelOldBackupAPPCommand(dir string) string {
	stCmd := "find -type f -name '*.zip' -mtime"
	ndCmd := "+" + strconv.Itoa(helpers.Conf.Backup.APP.Retain)
	rdCmd := "-delete"
	cmdSeries := []string{
		"cd " + dir,
		strings.Join([]string{stCmd, ndCmd, rdCmd}, " "),
	}
	return strings.Join(cmdSeries, ";")
}

// parseDelOldBackupDBCommand combine all args
func parseDelOldBackupDBCommand(dir string) string {
	stCmd := "find -type f -name '*.zip' -mtime"
	ndCmd := "+" + strconv.Itoa(helpers.Conf.Backup.DB.Retain)
	rdCmd := "-delete"
	cmdSeries := []string{
		"cd " + dir,
		strings.Join([]string{stCmd, ndCmd, rdCmd}, " "),
	}
	return strings.Join(cmdSeries, ";")
}

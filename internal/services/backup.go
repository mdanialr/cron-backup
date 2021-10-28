package services

import (
	"os"
	"sync"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
)

// Backup main function to control all backup that processed by goroutine
func Backup() {
	var wg sync.WaitGroup

	helpers.NzLogInfo.Println("Backup app invoked!")
	wg.Add(1)
	go backupAPP(&wg)

	helpers.NzLogInfo.Println("Backup database invoked!")
	wg.Add(1)
	go backupDB(&wg)

	wg.Wait()

	if err := deleteDumpedFile(); err != nil {
		helpers.NzLogError.Println(err)
	}
}

// makeSureDirExists make sure dir for backup apps exists by creating it
func makeSureDirExists(dir string) error {
	if err := os.MkdirAll(dir, 0770); err != nil {
		return err
	}
	return nil
}

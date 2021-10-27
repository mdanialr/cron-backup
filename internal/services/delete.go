package services

import (
	"os"
	"strings"
	"time"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
)

// deleteOldBackup delete old backup from the given dir according to
// its retain days
func deleteOldBackup(dir string, retainDay int) error {
	helpers.NzLogInfo.Println("[START] deleting old backup in:", "'"+dir+"'")
	currentDay := time.Now().Day()

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fl := range files {
		fInf, err := fl.Info()
		if err != nil {
			return err
		}

		day := fInf.ModTime().Day()
		if day < (currentDay - retainDay) {
			if err := os.Remove(dir + "/" + fl.Name()); err != nil {
				return err
			}
		}
	}
	helpers.NzLogInfo.Println("[DONE] deleting old backup in", "'"+dir+"'")
	return nil
}

// deleteDumpedFile delete dumped file in /tmp after zipping it
func deleteDumpedFile() error {
	helpers.NzLogInfo.Println("[START] deleting leftover dumped db file from /tmp")

	files, err := os.ReadDir("/tmp")
	if err != nil {
		return err
	}

	for _, fl := range files {
		if strings.HasPrefix(fl.Name(), "dump") {
			if err := os.Remove("/tmp/" + fl.Name()); err != nil {
				return err
			}
		}
	}
	helpers.NzLogInfo.Println("[DONE] deleting leftover dumped db file from /tmp")
	return nil
}

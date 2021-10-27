package services

import (
	"os"
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

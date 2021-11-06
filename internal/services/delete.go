package services

import (
	"math"
	"os"
	"strings"
	"time"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
)

// deleteOldBackup delete old backup from the given dir according to
// its retain days
func deleteOldBackup(dir string, retainDay int) error {
	helpers.NzLogInfo.Println("[START] deleting old backup in:", "'"+dir+"'")

	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, fl := range files {
		fInf, err := fl.Info()
		if err != nil {
			return err
		}

		sinceCreate := math.Round(time.Since(fInf.ModTime()).Hours())
		if int(sinceCreate) > ((retainDay * 24) - 1) {
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

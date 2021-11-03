package utils

import (
	"log"
	"os"
	"strings"
	"sync"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
)

// testDelete delete zip file that created by RunTest
func testDelete() {
	var wg sync.WaitGroup

	if helpers.TCond.IsDel {
		if !helpers.TCond.IsNoAPP {
			wg.Add(1)
			go testDeleteDir(&wg, testConf.BackupAppDir)
		}

		if !helpers.TCond.IsNoDB {
			wg.Add(1)
			go testDeleteDir(&wg, testConf.BackupDBDir)
		}

		wg.Add(1)
		go testDeleteDir(&wg, testConf.LogDir)
	} else {
		wg.Add(1)
		go loopAndDelete(&wg, fileToDelete.APPname)
		wg.Add(1)
		go loopAndDelete(&wg, fileToDelete.DBname)
	}

	// block this func until all delete process done
	wg.Wait()
}

// loopAndDelete loop through all strings which are should be
// the zip files, and delete all of them.
func loopAndDelete(wg *sync.WaitGroup, files []string) {
	defer wg.Done()
	for _, file := range files {
		wg.Add(1)
		go testDeleteZipFile(wg, file)
	}
}

// testDeleteDir delete dir and their contents recursively
func testDeleteDir(wg *sync.WaitGroup, dir string) {
	defer wg.Done()

	log.Println("[START] deleting test backup in", "'"+dir+"'")
	if err := os.RemoveAll(dir); err != nil {
		log.Println("[ERROR]", err)
	}
	log.Println("[DONE] deleting test backup in", "'"+dir+"'")
}

// testDeleteZipFile delete only zip file with given full path or dir
func testDeleteZipFile(wg *sync.WaitGroup, file string) {
	defer wg.Done()

	log.Println("[START] deleting zip file:", "'"+file+"'")
	if err := os.Remove(file); err != nil {
		log.Println("[ERROR]", err)
	}
	log.Println("[DONE] deleting zip file:", "'"+file+"'")
}

// testDeleteDumpedFile delete dumped file in /tmp after zipping it
func testDeleteDumpedFile() error {
	log.Println("[START] deleting leftover dumped db file from /tmp")

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
	log.Println("[DONE] deleting leftover dumped db file from /tmp")
	return nil
}

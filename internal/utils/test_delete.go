package utils

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
)

// testDelete delete zip file that created by RunTest
func testDelete() bool {
	isPass := true

	cAPP := make(chan bool)
	cDB := make(chan bool)
	cLOG := make(chan bool)

	if helpers.TCond.IsDel {
		if !helpers.TCond.IsNoAPP {
			go testDeleteDir(cAPP, testConf.BackupAppDir)
			if !<-cAPP {
				isPass = false
			}
		} else {
			go func() {
				cAPP <- true
			}()
		}

		if !helpers.TCond.IsNoDB {
			go testDeleteDir(cDB, testConf.BackupDBDir)
			if !<-cDB {
				isPass = false
			}
		} else {
			go func() {
				cDB <- true
			}()
		}

		go testDeleteDir(cLOG, testConf.LogDir)
		if !<-cLOG {
			isPass = false
		}
	} else {
		go testDeleteZipFile(cAPP, fileToDelete.APPname)
		go testDeleteZipFile(cDB, fileToDelete.DBname)
		go func() {
			cLOG <- true
		}()
	}
	if !isPass {
		os.Exit(1)
	}

	return isPass
}

// testDeleteDir delete dir and their contents recursively
func testDeleteDir(c chan bool, dir string) {
	isPass := true

	log.Println("[START] deleting test backup in", "'"+dir+"'")
	out, err := exec.Command("sh", "-c", "rm -r "+dir).CombinedOutput()
	if err != nil {
		log.Println(string(out))
		isPass = false
	}
	log.Println("[DONE] deleting test backup in", "'"+dir+"'")

	c <- isPass
}

// testDeleteZipFile delete only zip file with given full path or dir
func testDeleteZipFile(c chan bool, dir string) {
	isPass := true

	log.Println("[START] deleting zip file:", "'"+dir+"'")
	out, err := exec.Command("sh", "-c", "rm "+dir).CombinedOutput()
	if err != nil {
		log.Println(string(out))
		isPass = false
	}
	log.Println("[DONE] deleting zip file:", "'"+dir+"'")

	c <- isPass
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

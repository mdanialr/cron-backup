package utils

import (
	"log"
	"os/exec"
)

// testDelete delete zip file that created by RunTest
func testDelete(isDel bool) bool {
	isPass := true

	cAPP := make(chan bool)
	cDB := make(chan bool)
	cLOG := make(chan bool)

	if isDel {
		go testDeleteDir(cAPP, testConf.BackupAppDir)
		go testDeleteDir(cDB, testConf.BackupDBDir)
		go testDeleteDir(cLOG, testConf.LogDir)
	} else {
		go testDeleteZipFile(cAPP, fileToDelete.APPname)
		go testDeleteZipFile(cDB, fileToDelete.DBname)
		go func() {
			cLOG <- true
		}()
	}

	if !<-cAPP || !<-cDB || !<-cLOG {
		isPass = false
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

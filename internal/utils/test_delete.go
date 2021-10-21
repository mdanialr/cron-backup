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
		go testDeleteDBnAPPnLOG(cAPP, testConf.BackupAppDir)
		go testDeleteDBnAPPnLOG(cDB, testConf.BackupDBDir)
		go testDeleteDBnAPPnLOG(cLOG, testConf.LogDir)
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

// testDeleteDBnAPPnLOG delete dir and their contents recursively
func testDeleteDBnAPPnLOG(c chan bool, dir string) {
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

	log.Println("[START] deleting test backup in", "'"+dir+"'")
	out, err := exec.Command("sh", "-c", "rm "+dir).CombinedOutput()
	if err != nil {
		log.Println(string(out))
		isPass = false
	}
	log.Println("[DONE] deleting test backup in", "'"+dir+"'")

	c <- isPass
}

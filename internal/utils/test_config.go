package utils

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// testCheckConfig check some important params must not empty
func testCheckConfig() bool {
	isPass := true

	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		errMsg := "[ERROR] failed to load yaml file; " + err.Error()
		log.Fatalln(errMsg)
		return false
	}

	if err = yaml.Unmarshal(yamlFile, &testConf); err != nil {
		errMsg := "[ERROR] Unmarshal: " + err.Error()
		log.Fatalln(errMsg)
		return false
	}

	if testConf.LogDir == "" {
		log.Println("[ERROR] Log dir {log_dir:} in config.yaml file must not empty!")
		isPass = false
	}
	if testConf.Backup.RootDir == "" {
		log.Println("[ERROR] Backup root dir {root_dir:} in config.yaml file must not empty!")
		isPass = false
	}
	if testConf.Backup.Retain == 0 {
		log.Println("[ERROR] Number of days to retain backup (app and db) in config.yaml file must not empty!")
		isPass = false
	}
	if len(testConf.Backup.APP.Apps) == 0 {
		log.Println("[ERROR] There should be at least one appfile {- appfile:} configured!")
		isPass = false
	}
	if len(testConf.Backup.DB.Databases) == 0 {
		log.Println("[ERROR] There should be at least one database {- database:} configured!")
		isPass = false
	}
	if err := testConf.EnsureDBTypeExists(); err != nil {
		log.Println("[ERROR] Make sure DB Type is not empty. Fill in with either 'pg' or 'mdb'.")
		isPass = false
	}
	testConf.SetupDBType()
	if err := testConf.SanitizeAndCheckDB(); err != nil {
		log.Println("[ERROR]", err)
		isPass = false
	}
	if !isPass {
		os.Exit(1)
	}

	testConf.SanitizeLogDir()
	testConf.SanitizeRootDir()
	testConf.SanitizeAppDir()
	testConf.SetupBackupDir()
	testConf.SetupSpecificBackupRetain()
	testConf.SanitizeAndSetupSample()

	return isPass
}

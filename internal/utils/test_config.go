package utils

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
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

	AssignSampleFromFlag()
	testConf.SanitizeLogDir()
	testConf.SanitizeRootDir()
	testConf.SanitizeAppDir()
	testConf.SetupBackupDir()
	testConf.SetupSpecificBackupRetain()
	testConf.SanitizeAndSetupSample()

	return isPass
}

// AssignSampleFromFlag assign config sample's value from
// cli args (flag).
func AssignSampleFromFlag() {
	if !isFlagExist("sam-app") {
		helpers.TCond.Sapp = helpers.TCond.Sample
	}
	if !isFlagExist("sam-db") {
		helpers.TCond.Sdb = helpers.TCond.Sample
	}

	testConf.Backup.APP.Sample = helpers.TCond.Sapp
	testConf.Backup.DB.Sample = helpers.TCond.Sdb
}

// isFlagExist check if a flag is set or not.
func isFlagExist(name string) bool {
	var founded bool
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			founded = true
		}
	})

	return founded
}

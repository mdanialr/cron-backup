package helpers

import (
	"io/ioutil"
	"log"

	"github.com/mdanialr/go-cron-backup/internal/models"
	"gopkg.in/yaml.v2"
)

var Conf *models.Config

// LoadConfigFromFile load config.yaml file and assign it to Config
func LoadConfigFromFile() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		errMsg := "Error when load yaml file; " + err.Error()
		log.Fatalln(errMsg)
	}

	if err = yaml.Unmarshal(yamlFile, &Conf); err != nil {
		errMsg := "Unmarshal: " + err.Error()
		log.Fatalln(errMsg)
	}

	Conf.SanitizeLogDir()
	Conf.SanitizeRootDir()
	Conf.SanitizeAppDir()
	Conf.SetupBackupDir()
	Conf.SetupSpecificBackupRetain()
	Conf.EnsureDBTypeExists()
	Conf.SetupDBType()
	Conf.SanitizeAndCheckDB()
}

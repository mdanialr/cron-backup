package arch

import (
	"os/exec"
	"strings"
	"time"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
	"github.com/mdanialr/go-cron-backup/internal/models"
)

// BashDBZip create the fullpath plus formated zip name from
// dumped database
func BashDBZip(db models.Database, fileToZip string) error {
	// create the z
	fmtTime := time.Now().Format("2006-Jan-02_Monday_15:04:05")
	fName := "/" + fmtTime + ".zip"
	zipName := helpers.Conf.BackupDBDir + db.DirName + fName

	if err := zipDBWithBash(zipName, fileToZip); err != nil {
		return err
	}
	return nil
}

// zipDBWithBash zip the given zipName from fileToZip file
func zipDBWithBash(zipName string, fileToZip string) error {
	cmdSeries := []string{
		"cd /tmp",
		"zip -q " + zipName + " " + fileToZip,
	}
	cmd := strings.Join(cmdSeries, ";")

	if _, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

// BashAPPZip create the fullpath plus formated zip name from
// directory target to zip
func BashAPPZip(app models.App) error {
	fmtTime := time.Now().Format("2006-Jan-02_Monday_15:04:05")
	fName := "/" + fmtTime + ".zip"
	zipName := helpers.Conf.BackupAppDir + app.DirName + fName

	if err := zipAPPWithBash(zipName, app); err != nil {
		return err
	}
	return nil
}

// zipAPPWithBash zip the given zipName from models.App value
func zipAPPWithBash(zipName string, app models.App) error {
	cmdSeries := []string{
		"cd " + app.AppDir,
		"zip -r -q " + zipName + " *",
	}
	cmd := strings.Join(cmdSeries, ";")

	if _, err := exec.Command("sh", "-c", cmd).CombinedOutput(); err != nil {
		return err
	}
	return nil
}

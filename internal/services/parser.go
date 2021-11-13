package services

import (
	"strconv"
	"strings"

	"github.com/mdanialr/go-cron-backup/internal/models"
)

// parseDumpingMariaDBCommand combine all commands for dumping database
func parseDumpingMariaDBCommand(db models.Database) (string, string) {
	cmd := "mariadb-dump " + db.Name
	host := "-h " + db.Host
	port := "-P " + strconv.Itoa(db.Port)
	usr := "-u " + db.Usr
	pwd := "-p" + db.Pwd
	opt_params := db.OptParams
	outName := "dump_" + db.Name
	cmdSeries := []string{
		cmd,
		host,
		port,
		usr,
		pwd,
		opt_params,
		">",
		outName,
	}
	dumpCmd := strings.Join(cmdSeries, " ")
	return strings.Join([]string{"cd /tmp", dumpCmd}, ";"), outName
}

// parseDumpingPGCommand combine all commands for dumping database
func parseDumpingPGCommand(db models.Database) (string, string) {
	cmd := "pg_dump "
	params := "--dbname=postgresql://"
	creds := db.Usr + ":" + db.Pwd + "@"
	sock := db.Host + ":" + strconv.Itoa(db.Port) + "/"
	outName := "dump_" + db.Name
	cmdSeries := []string{
		cmd,
		params,
		creds,
		sock,
		db.Name,
	}
	dumpCmd := strings.Join(cmdSeries, "")
	dumpCmd += strings.Join([]string{">", outName}, " ")

	return strings.Join([]string{"cd /tmp", dumpCmd}, ";"), outName
}

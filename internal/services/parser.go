package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mdanialr/go-cron-backup/internal/models"
)

// parseDumpingMysqlCommand combine all commands for dumping database
func parseDumpingMysqlCommand(db models.Database) (string, string) {
	cmd := fmt.Sprintf("mysqldump -h %s -P %d -u %s -p%s %s %s", db.Host, db.Port, db.Usr, db.Pwd, db.Name, db.OptParams)
	outName := fmt.Sprintf("dump_%s_%s", db.Name, time.Now().Format("2006-01-02_15-04-05"))
	cmd = fmt.Sprintf("%s > %s", cmd, outName)
	return fmt.Sprintf("cd /tmp; %s", cmd), outName
}

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

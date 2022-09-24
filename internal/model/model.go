package model

import (
	"fmt"
	"log"
	"strings"
)

type (
	// App detail information about an app/directory to be backed up.
	App struct {
		Dir      string `mapstructure:"dir"`  // target directory that will be archived.
		Name     string `mapstructure:"name"` // directory name where this backup is stored.
		StoreDir string `mapstructure:"-"`    // directory where this backup is stored after got archived.
	}
	// APP the most outer struct of the config file containing all config info for apps/directories.
	APP struct {
		MaxDays   uint   `mapstructure:"max_days"`
		MaxWorker uint   `mapstructure:"max_worker"`
		Apps      []*App `mapstructure:"apps"`
	}

	// Database detail information about a database to be backed up.
	Database struct {
		ID         string `mapstructure:"-"`      // unique id for the database to prevent collision with another database when dumping it.
		Type       string `mapstructure:"type"`   // either my: mysql, pg: postgres or md: mariadb
		Docker     string `mapstructure:"docker"` // docker container name
		Host       string `mapstructure:"host"`   // hostname or ip address. if docker is set, this is ignored
		Port       uint   `mapstructure:"port"`   // port number. if docker is set, this is ignored
		Name       string `mapstructure:"name"`   // database name
		User       string `mapstructure:"user"`   // database user
		Pass       string `mapstructure:"pass"`   // database password
		BackupName string `mapstructure:"backup"` // directory name where this backup is stored
		OptParams  string `mapstructure:"params"` // optional parameters
		CMD        string `mapstructure:"-"`      // parsed command to be executed
		Dir        string `mapstructure:"-"`      // directory where the backup is stored
	}
	// DB the most outer struct of the config file containing all config info for databases.
	DB struct {
		MaxDays   uint        `mapstructure:"max_days"`
		MaxWorker uint        `mapstructure:"max_worker"`
		Databases []*Database `mapstructure:"databases"`
	}

	// Logs bag for both info and error logger.
	Logs struct {
		Inf *log.Logger // Inf is the logger for info messages
		Err *log.Logger // Err is the logger for error messages
	}
)

// SetStoreDir append the given root with the app name without adding any trailing slash then assign it to StoreDir
// field.
func (a *App) SetStoreDir(root string) {
	a.StoreDir = fmt.Sprintf("%s/%s", root, a.Name)
}

// setType normalize and set the type of the database.
func (d *Database) setType() error {
	if strings.HasPrefix(d.Type, "my") {
		d.Type = "my"
		return nil
	}
	if strings.HasPrefix(d.Type, "pg") {
		d.Type = "pg"
		return nil
	}
	if strings.HasPrefix(d.Type, "md") {
		d.Type = "md"
		return nil
	}

	return fmt.Errorf("unsupported database type: '%s'. currently supported types are [pg,my,md]", d.Type)
}

// setDefault give default value to the Database instance's fields such as Host, Port & User.
func (d *Database) setDefault() {
	if d.Host == "" {
		d.Host = "localhost"
	}
	if d.Port == 0 {
		switch d.Type {
		case "my", "md":
			d.Port = 3306
		case "pg":
			d.Port = 5432
		}
	}
	if d.User == "" {
		switch d.Type {
		case "my", "md":
			d.User = "root"
		case "pg":
			d.User = "postgres"
		}
	}
}

// buildID build unique id for the database.
func (d *Database) buildID() {
	d.ID = fmt.Sprintf("%s_%s_%d_%s", d.Type, d.Host, d.Port, d.Name)
}

// buildCMD build the command to be executed for this Database instance.
func (d *Database) buildCMD() string {
	cmd := fmt.Sprintf("-h %s -P %d -u%s -p%s %s %s", d.Host, d.Port, d.User, d.Pass, d.OptParams, d.Name)
	if len(d.Pass) < 1 {
		cmd = strings.ReplaceAll(cmd, "-p", "")
	}
	switch d.Type {
	case "my":
		cmd = fmt.Sprintf("mysqldump %s", cmd)
	case "md":
		cmd = fmt.Sprintf("mariadb-dump %s", cmd)
	case "pg":
		cmd = fmt.Sprintf("pg_dump postgresql://%s:%s@%s:%d/%s", d.User, d.Pass, d.Host, d.Port, d.Name)
	}

	if d.Docker != "" {
		return fmt.Sprintf("docker exec -t %s %s", d.Docker, cmd)
	}

	return cmd
}

// SetDir append the given root with the backup name without adding any trailing slash then assign it to Dir field.
func (d *Database) SetDir(root string) {
	d.Dir = fmt.Sprintf("%s/%s", root, d.BackupName)
}

// Config a bag for both DB and APP instance.
type Config struct {
	DB
	APP
}

// Validate do various validation rules. return error if any required field is not provided or invalid.
func (c *Config) Validate() error {
	if len(c.DB.Databases) > 0 {
		if err := c.checkDuplicateDB(c.DB.Databases); err != nil {
			return fmt.Errorf("found duplicate database backup name: %s", err)
		}
		for _, db := range c.DB.Databases {
			if err := db.setType(); err != nil {
				return err
			}
			db.setDefault()
			db.buildID()
			db.CMD = db.buildCMD()
		}
	}

	if len(c.APP.Apps) > 0 {
		if err := c.checkDuplicateAPP(c.APP.Apps); err != nil {
			return fmt.Errorf("found duplicate app backup name: %s", err)
		}
	}

	return nil
}

// checkDuplicateDB return error if there is any duplicate in database's backup name.
func (c *Config) checkDuplicateDB(dbs []*Database) error {
	keys := make(map[string]bool)
	for _, entry := range dbs {
		if _, value := keys[entry.BackupName]; value {
			return fmt.Errorf("duplicate backup name (%s)", entry.BackupName)
		}
		keys[entry.BackupName] = true
	}

	return nil
}

// checkDuplicateAPP return error if there is any duplicate in app's name.
func (c *Config) checkDuplicateAPP(apps []*App) error {
	keys := make(map[string]bool)
	for _, entry := range apps {
		if _, value := keys[entry.Name]; value {
			return fmt.Errorf("duplicate app name (%s)", entry.Name)
		}
		keys[entry.Name] = true
	}

	return nil
}

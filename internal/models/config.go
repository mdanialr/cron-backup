package models

import (
	"errors"
	"strings"
)

type App struct {
	AppDir  string `yaml:"app_dir"`
	DirName string `yaml:"backup_dir_name"`
}

type Apps []struct {
	App App `yaml:"appfile"`
}

type APP struct {
	Apps      Apps `yaml:"apps"`
	Retain    int  `yaml:"max_days_to_retain"`
	MaxWorker int  `yaml:"max_worker"`
	Sample    int
}

type dbType struct {
	PGsql   bool
	MariaDB bool
}

type Database struct {
	T         dbType
	Type      string `yaml:"type"`
	Host      string `yaml:"hostname"`
	Port      int    `yaml:"port"`
	Name      string `yaml:"dbname"`
	Usr       string `yaml:"dbuser"`
	Pwd       string `yaml:"dbpass"`
	DirName   string `yaml:"backup_dir_name"`
	OptParams string `yaml:"opt_params"`
}

type Databases []struct {
	Database Database `yaml:"database"`
}

type DB struct {
	Databases Databases `yaml:"databases"`
	Retain    int       `yaml:"max_days_to_retain"`
	MaxWorker int       `yaml:"max_worker"`
	Sample    int
}

type Backup struct {
	DB      DB     `yaml:"database"`
	APP     APP    `yaml:"app"`
	Retain  int    `yaml:"max_days_to_retain"`
	RootDir string `yaml:"root_dir"`
}

type Config struct {
	LogDir       string `yaml:"log_dir"`
	Backup       Backup `yaml:"backup"`
	BackupDBDir  string
	BackupAppDir string
}

// SanitizeAppDir make sure all app dir has trailing slash
func (c *Config) SanitizeAppDir() {
	for i := range c.Backup.APP.Apps {
		v := &c.Backup.APP.Apps[i]
		if !strings.HasSuffix(v.App.AppDir, "/") {
			v.App.AppDir = v.App.AppDir + "/"
		}
	}
}

// SanitizeRootDir make sure backup root dir has trailing slash
func (c *Config) SanitizeRootDir() {
	if !strings.HasSuffix(c.Backup.RootDir, "/") {
		c.Backup.RootDir = c.Backup.RootDir + "/"
	}
}

// SanitizeLogDir make sure log dir has trailing slash
func (c *Config) SanitizeLogDir() {
	if !strings.HasSuffix(c.LogDir, "/") {
		c.LogDir = c.LogDir + "/"
	}
}

// SetupBackupDir add new dir to each backup type {db|app}
func (c *Config) SetupBackupDir() {
	c.BackupAppDir = c.Backup.RootDir + "app/"
	c.BackupDBDir = c.Backup.RootDir + "db/"
}

// SetupSpecificBackupRetain assign global retain if specific
// retain not specified
func (c *Config) SetupSpecificBackupRetain() {
	if c.Backup.DB.Retain == 0 {
		c.Backup.DB.Retain = c.Backup.Retain
	}
	if c.Backup.APP.Retain == 0 {
		c.Backup.APP.Retain = c.Backup.Retain
	}
}

// EnsureDBTypeExists return error if DB type is empty
func (c *Config) EnsureDBTypeExists() error {
	for i := range c.Backup.DB.Databases {
		v := &c.Backup.DB.Databases[i]
		if v.Database.Type == "" {
			return errors.New("")
		}
	}
	return nil
}

// SanitizeAndCheckDB return error if DB user is empty.
// Then assign default value to hostname and port if
// either are empty.
func (c *Config) SanitizeAndCheckDB() error {
	for i := range c.Backup.DB.Databases {
		db := &c.Backup.DB.Databases[i].Database
		if db.Usr == "" {
			return errors.New("make sure user for database connection is not empty")
		}
		if db.Host == "" {
			db.Host = "localhost"
		}
		if db.Port == 0 {
			if db.T.MariaDB {
				db.Port = 3306
			}
			if db.T.PGsql {
				db.Port = 5432
			}
		}
	}
	return nil
}

// SetupDBType distinguish PostgresSQL and MariaDB and setup
// the bool type
func (c *Config) SetupDBType() {
	for i := range c.Backup.DB.Databases {
		v := &c.Backup.DB.Databases[i]
		l := strings.ToLower(v.Database.Type)
		if strings.HasPrefix(l, "pg") {
			v.Database.T.PGsql = true
		}
		if strings.HasPrefix(l, "my") || strings.HasPrefix(l, "md") {
			v.Database.T.MariaDB = true
		}
	}
}

// SanitizeAndSetupSample sanitize to prevent panic caused slice
// bounds out of range or sample is zero or not set
func (c *Config) SanitizeAndSetupSample() {
	if c.Backup.APP.Sample == 0 {
		c.Backup.APP.Sample = 1
	}
	if c.Backup.DB.Sample == 0 {
		c.Backup.DB.Sample = 1
	}
	if c.Backup.APP.Sample > len(c.Backup.APP.Apps) {
		c.Backup.APP.Sample = len(c.Backup.APP.Apps)
	}
	if c.Backup.DB.Sample > len(c.Backup.DB.Databases) {
		c.Backup.DB.Sample = len(c.Backup.DB.Databases)
	}
}

// SanitizeMaxWorker give default value to empty max_worker.
func (c *Config) SanitizeMaxWorker() {
	if c.Backup.APP.MaxWorker == 0 {
		c.Backup.APP.MaxWorker = 1
	}
	if c.Backup.DB.MaxWorker == 0 {
		c.Backup.DB.MaxWorker = 1
	}
}

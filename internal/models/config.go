package models

import (
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
	Apps   Apps `yaml:"apps"`
	Retain int  `yaml:"days_number_to_retain"`
}

type dbType struct {
	PGsql   bool
	MariaDB bool
}

type Database struct {
	T       dbType
	Type    string `yaml:"type"`
	Name    string `yaml:"dbname"`
	Usr     string `yaml:"dbuser"`
	Pwd     string `yaml:"dbpass"`
	DirName string `yaml:"backup_dir_name"`
}

type Databases []struct {
	Database Database `yaml:"database"`
}

type DB struct {
	Databases Databases `yaml:"databases"`
	Retain    int       `yaml:"days_number_to_retain"`
}

type Backup struct {
	DB      DB     `yaml:"database"`
	APP     APP    `yaml:"app"`
	Retain  int    `yaml:"days_number_to_retain"`
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

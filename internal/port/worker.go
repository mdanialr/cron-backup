package port

import (
	"bytes"

	"github.com/mdanialr/go-cron-backup/internal/model"
)

type (
	DBJob interface {
		// Dump a worker for DB that run the CMD (dumping db) for the given database instance and return the result in buffer.
		Dump(db *model.Database) (*bytes.Buffer, error)
	}

	DeleteJob interface {
		// DeleteOldBackup delete all files in dir based on the given retention days.
		DeleteOldBackup(dir string, retain uint) error
	}
)

package worker

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"

	"github.com/mdanialr/go-cron-backup/internal/model"
	"github.com/mdanialr/go-cron-backup/internal/port"
	"github.com/mdanialr/go-cron-backup/pkg/helper"
)

// NewDBJob create concrete implementation of port.DBJob.
func NewDBJob(os string, log *model.Logs) port.DBJob {
	return &dbJob{
		log: log,
		os:  os,
	}
}

type dbJob struct {
	log *model.Logs
	os  string
}

func (w *dbJob) Dump(db *model.Database) (*bytes.Buffer, error) {
	defer w.log.Inf.Printf(helper.LogDone(db.ID, "dumping database"))
	w.log.Inf.Printf(helper.LogStart(db.ID, "dumping database"))

	cmdName, cmdArg := "sh", "-c" // for linux and darwin (macOS)
	if w.os == "windows" {
		cmdName, cmdArg = "cmd", "/C"
	}

	stdOut, stdErr := new(bytes.Buffer), new(bytes.Buffer)
	cmd := exec.Command(cmdName, cmdArg, db.CMD)
	cmd.Stdout = stdOut
	cmd.Stderr = stdErr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run command for dumping database ( %s ): %s", db.ID, err)
	}
	if stdErr.Len() > 0 {
		return nil, fmt.Errorf("got error in stderr: %s", stdErr.String())
	}

	return stdOut, nil
}

// NewDeleteJob create concrete implementation of port.DeleteJob.
func NewDeleteJob(log *model.Logs) port.DeleteJob {
	return &deleteJob{
		log: log,
	}
}

type deleteJob struct {
	log *model.Logs
}

// DeleteOldBackup delete old backup based on the given retention days.
func (a *deleteJob) DeleteOldBackup(dir string, retain uint) error {
	defer a.log.Inf.Printf(helper.LogDone(dir, "deleting old backup"))
	a.log.Inf.Printf(helper.LogStart(dir, "deleting old backup"))

	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read dir: %s", err)
	}

	for _, fl := range files {
		inf, err := fl.Info()
		if err != nil {
			return fmt.Errorf("failed to get file info: %s", err)
		}

		sinceCreated := math.Round(time.Since(inf.ModTime()).Hours())
		if uint(sinceCreated) > ((retain * 24) - 1) {
			if err = os.Remove(fmt.Sprintf("%s/%s", dir, fl.Name())); err != nil {
				return fmt.Errorf("failed to delete file: %s", err)
			}
		}
	}

	return nil
}

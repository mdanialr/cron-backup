package worker

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
	"github.com/mdanialr/go-cron-backup/internal/model"
	"github.com/mdanialr/go-cron-backup/internal/port"
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
	defer w.log.Inf.Printf(helpers.LogDone(db.ID, "dumping database"))
	w.log.Inf.Printf(helpers.LogStart(db.ID, "dumping database"))

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

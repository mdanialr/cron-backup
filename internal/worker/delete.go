package worker

import (
	"sync"

	"github.com/mdanialr/go-cron-backup/internal/model"
)

// DeleteWorker worker that will do the job which is deleting old backup files.
func DeleteWorker(wg *sync.WaitGroup, jobs <-chan *model.DeleteWorker, log *model.Logs) {
	for job := range jobs {
		newJob := NewDeleteJob(log)
		if err := newJob.DeleteOldBackup(job.Dir, job.Retain); err != nil {
			log.Err.Println(err)
			wg.Done()
			continue
		}

		wg.Done()
	}
}

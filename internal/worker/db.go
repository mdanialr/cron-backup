package worker

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/mdanialr/go-cron-backup/internal/model"
	"github.com/mdanialr/go-cron-backup/pkg/archive"
	"github.com/mdanialr/go-cron-backup/pkg/helper"
)

// DBWorker worker that will do the job which is dumping the database, zipping it, then write it to a file.
func DBWorker(wg *sync.WaitGroup, jobs <-chan *model.Database, goos string, log *model.Logs) {
	for db := range jobs {
		newJob := NewDBJob(goos, log)
		dumped, err := newJob.Dump(db)
		if err != nil {
			log.Err.Println(err)
			wg.Done()
			continue
		}

		// zip the buffer
		log.Inf.Printf(helper.LogStart(db.ID, "zipping"))
		zipped, err := archive.ZipFile(fmt.Sprintf("%s_dump", db.Name), dumped.Bytes())
		if err != nil {
			log.Err.Println("failed to zip dumped database:", err)
			log.Inf.Printf(helper.LogDone(db.ID, "zipping"))
			wg.Done()
			continue
		}
		log.Inf.Printf(helper.LogDone(db.ID, "zipping"))

		// write to file
		timeNow := time.Now().Format("2006-Jan-02_15-04-05")
		log.Inf.Printf(helper.LogStart(db.ID, "writing to file"))
		fl, err := os.Create(fmt.Sprintf("%s/%s", db.Dir, timeNow))
		if err != nil {
			log.Err.Println("failed to create file:", err)
			log.Inf.Printf(helper.LogDone(db.ID, "writing to file"))
			wg.Done()
			continue
		}

		fl.Write(zipped.Bytes())
		fl.Close()
		log.Inf.Printf(helper.LogDone(db.ID, "writing to file"))
		wg.Done()
	}
}

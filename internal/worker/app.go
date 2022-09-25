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

// APPWorker worker that will zip the given directory and write it to a file via jobs channel.
func APPWorker(wg *sync.WaitGroup, jobs <-chan *model.App, log *model.Logs) {
	for app := range jobs {
		// zip the directory
		log.Inf.Printf(helper.LogStart(app.Name, "zipping"))
		zipped, err := archive.ZipDir(app.Dir)
		if err != nil {
			log.Err.Println("failed to zip directory:", err)
			log.Inf.Printf(helper.LogDone(app.Name, "zipping"))
			wg.Done()
			continue
		}
		log.Inf.Printf(helper.LogDone(app.Name, "zipping"))

		// write to file
		timeNow := time.Now().Format("2006-Jan-02_15-04-05")
		log.Inf.Printf(helper.LogStart(app.Name, "writing to file"))
		fl, err := os.Create(fmt.Sprintf("%s/%s.zip", app.StoreDir, timeNow))
		if err != nil {
			log.Err.Println("failed to create file:", err)
			log.Inf.Printf(helper.LogDone(app.Name, "writing to file"))
			wg.Done()
			continue
		}

		fl.Write(zipped.Bytes())
		fl.Close()
		log.Inf.Printf(helper.LogDone(app.Name, "writing to file"))
		wg.Done()
	}
}

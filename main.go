package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/mdanialr/go-cron-backup/internal/model"
	"github.com/mdanialr/go-cron-backup/internal/worker"
	"github.com/mdanialr/go-cron-backup/pkg/config"
	"github.com/mdanialr/go-cron-backup/pkg/logger"
)

func main() {
	conf, err := config.InitConfig(".")
	if err != nil {
		log.Fatalln("failed to init config file:", err)
	}
	if err = config.SetupDefault(conf); err != nil {
		log.Fatalln("failed to sanitize and setup default config:", err)
	}

	var svc model.Config
	if err = conf.Unmarshal(&svc); err != nil {
		log.Fatalln("failed to unmarshal config:", err)
	}
	if err = svc.Validate(); err != nil {
		log.Fatalln("failed to validate config file:", err)
	}

	infLog, err := logger.InitInfoLogger(conf)
	if err != nil {
		log.Fatalln("failed to init info logger:", err)
	}
	errLog, err := logger.InitErrorLogger(conf)
	if err != nil {
		log.Fatalln("failed to init error logger:", err)
	}
	logBag := &model.Logs{Inf: infLog, Err: errLog}

	// make sure target dir is accessible and writable
	for _, db := range svc.DB.Databases {
		db.SetDir(fmt.Sprintf("%s/databases", conf.GetString("root")))
		if err = os.MkdirAll(db.Dir, 0770); err != nil {
			logBag.Err.Println("failed to create dir:", err)
			return
		}
	}
	for _, app := range svc.APP.Apps {
		app.SetStoreDir(fmt.Sprintf("%s/apps", conf.GetString("root")))
		if err = os.MkdirAll(app.StoreDir, 0770); err != nil {
			logBag.Err.Println("failed to create dir:", err)
			return
		}
	}

	var wg sync.WaitGroup
	dbJobChan := make(chan *model.Database)
	appJobChan := make(chan *model.App)

	// spawn as many workers as many in the config
	for i := 1; i <= int(svc.DB.MaxWorker); i++ {
		go worker.DBWorker(&wg, dbJobChan, runtime.GOOS, logBag)
	}
	for i := 1; i <= int(svc.APP.MaxWorker); i++ {
		go worker.APPWorker(&wg, appJobChan, logBag)
	}

	logBag.Inf.Println("")
	logBag.Inf.Println("-------------- START --------------")

	// send jobs
	for _, job := range svc.DB.Databases {
		wg.Add(1)
		dbJobChan <- job
	}
	close(dbJobChan)
	for _, job := range svc.APP.Apps {
		wg.Add(1)
		appJobChan <- job
	}
	close(appJobChan)

	wg.Wait()
	logBag.Inf.Println("-------------- DONE --------------")
	logBag.Inf.Println("")
}

package main

import (
	"flag"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
	"github.com/mdanialr/go-cron-backup/internal/services"
	"github.com/mdanialr/go-cron-backup/internal/utils"
)

func main() {
	setupFlag()

	// Run backup only if there is no test flag
	if !helpers.TCond.IsTest {
		// Initialize Config & Logger
		helpers.LoadConfigFromFile()
		helpers.InitNzLog()
		// Start backup app & database process
		helpers.NzLogInfo.Println("")
		helpers.NzLogInfo.Println("=== Invoking backup ===")
		services.Backup()
		helpers.NzLogInfo.Println("=== Backup successfully invoked ===")
		helpers.NzLogInfo.Println("")
		return
	}
	utils.RunTest()
}

// setupFlag initialize flag then parse it.
func setupFlag() {
	flag.BoolVar(&helpers.TCond.IsTest, "test", false, "to test this app functioning properly")
	flag.BoolVar(&helpers.TCond.IsDel, "d", false, "to delete all created dir after this test")
	flag.BoolVar(&helpers.TCond.IsNoDB, "no-db", false, "to exclude database from this backup test")
	flag.BoolVar(&helpers.TCond.IsNoAPP, "no-app", false, "to exclude app from this backup test")
	flag.IntVar(&helpers.TCond.Sample, "sample", 1, "set number of samples for both app and database")
	flag.IntVar(&helpers.TCond.Sapp, "sam-app", 1, "spesifically set number of samples for app")
	flag.IntVar(&helpers.TCond.Sdb, "sam-db", 1, "spesifically set number of samples for database")
	flag.Parse()
}

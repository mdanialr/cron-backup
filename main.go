package main

import (
	"flag"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
	"github.com/mdanialr/go-cron-backup/internal/services"
	"github.com/mdanialr/go-cron-backup/internal/utils"
)

func main() {
	isTest := flag.Bool("test", false, "to test this app functioning properly")
	isDel := flag.Bool("d", false, "to delete all created dir after this test")
	isExDB := flag.Bool("no-db", false, "to exclude database from this backup test")
	isExAPP := flag.Bool("no-app", false, "to exclude app from this backup test")
	flag.Parse()

	// Run backup only if there is no test flag
	if !*isTest {
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
	utils.RunTest(*isDel, *isExDB, *isExAPP)
}

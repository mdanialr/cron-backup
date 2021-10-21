package main

import (
	"flag"

	"github.com/mdanialr/go-cron-backup/internal/helpers"
	"github.com/mdanialr/go-cron-backup/internal/services"
	"github.com/mdanialr/go-cron-backup/internal/utils"
)

func main() {
	test := flag.Bool("test", false, "to test this app functioning properly")
	flag.Parse()

	// Run backup only if there is no test flag
	if !*test {
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

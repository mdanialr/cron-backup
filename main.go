package main

import (
	"github.com/mdanialr/go-cron-backup/internal/helpers"
	"github.com/mdanialr/go-cron-backup/internal/services"
)

func main() {
	// Initialize Logger
	helpers.LoadConfigFromFile()
	helpers.InitNzLog()

	// Start backup app & database process
	helpers.NzLogInfo.Println("")
	helpers.NzLogInfo.Println("=== Invoking backup ===")
	services.Backup()
	helpers.NzLogInfo.Println("=== Backup successfully invoked ===")
	helpers.NzLogInfo.Println("")
}

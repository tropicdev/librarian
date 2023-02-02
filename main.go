package main

import (
	"bitecodelabs.com/librarian/config"
	"bitecodelabs.com/librarian/logger"
	"github.com/robfig/cron"
)

func main() {

	logger.InfoLog.Println("Librarian is starting")
	config := config.LoadConfig()

	c := cron.New()
	for _, server := range config.Pterodactyl.Servers {
		err := c.AddFunc(server.Schedule, func() {
			logger.InfoLog.Println("Getting Backup for", server.Name)
			// run the command here
		})
		if err != nil {
			logger.ErrorLog.Fatalf("error scheduling backup: %v", err)
		}
	}
	c.Start()

}

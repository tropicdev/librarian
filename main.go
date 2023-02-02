package main

import (
	"fmt"

	"bitecodelabs.com/librarian/backup"
	"bitecodelabs.com/librarian/config"
	"bitecodelabs.com/librarian/logger"
	"github.com/robfig/cron/v3"
)

func main() {

	logger.InfoLog.Println("Librarian is starting")
	json_config := config.LoadConfig()

	c := cron.New()
	for _, server := range json_config.Pterodactyl.Servers {
		fmt.Println(json_config)
		fmt.Println(server)
		config := config.BackupConfig{
			Server_Id:        server.ID,
			Host:             json_config.Pterodactyl.Host,
			API_Token:        json_config.Pterodactyl.APIToken,
			Name:             server.Name,
			Output_Directory: server.OutputDirectory,
		}
		_, err := c.AddFunc(server.Schedule, func() {
			logger.InfoLog.Println("Getting Backup for", server.Name)
			backup.Backup(config)
		})
		if err != nil {
			logger.ErrorLog.Fatalf("error scheduling backup: %v", err)
		}
	}
	c.Start()

	select {}
}

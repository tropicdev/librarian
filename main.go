package main

import (
	"os"
	"path/filepath"
	"time"

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
		config := config.BackupConfig{
			Server_Id:        server.ID,
			Host:             json_config.Pterodactyl.Host,
			API_Token:        json_config.Pterodactyl.APIToken,
			Name:             server.Name,
			Output_Directory: server.OutputDirectory,
			Delete_After:     server.Delete_After,
		}
		_, err := c.AddFunc(server.Schedule, func() {
			logger.InfoLog.Println("Getting Backup for", server.Name)
			backup.Backup(config)
		})
		if err != nil {
			logger.ErrorLog.Fatalf("error scheduling backup: %v", err)
		}

		if config.Delete_After > 0 {
			_, err = c.AddFunc("0 0 * * *", func() {
				dir := config.Output_Directory
				now := time.Now()
				ageLimit := time.Hour * 24 * time.Duration(config.Delete_After) // files older than X Days
				filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
					if info.IsDir() {
						return nil
					}
					if now.Sub(info.ModTime()) > ageLimit {
						logger.InfoLog.Printf("Deleting old file: %s\n", path)
						err := os.Remove(path)
						if err != nil {
							logger.ErrorLog.Printf("Error deleting file: %v\n", err)
						}
					}
					return nil
				})
			})
			if err != nil {
				logger.ErrorLog.Fatalf("error scheduling delete job: %v", err)
			}
		}

	}
	c.Start()

	select {}
}

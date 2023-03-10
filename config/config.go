package config

import (
	"encoding/json"
	"os"

	"bitecodelabs.com/librarian/logger"
)

type Config struct {
	Pterodactyl struct {
		Enabled  bool   `json:"enabled"`
		Host     string `json:"host"`
		APIToken string `json:"api_token"`
		Servers  []struct {
			ID              string `json:"id"`
			Name            string `json:"name"`
			OutputDirectory string `json:"output_directory"`
			Schedule        string `json:"schedule"`
			Delete_After    int    `json:"delete_after"`
		} `json:"servers"`
	} `json:"pterodactyl"`
	Database struct {
		Enabled   bool `json:"enabled"`
		Databases []struct {
			Host             string `json:"host"`
			User             string `json:"user"`
			Password         string `json:"password"`
			DBName           string `json:"db_name"`
			Schedule         string `json:"schedule"`
			Output_Directory string `json:"output_directory"`
			Delete_After     int    `json:"delete_after"`
		} `json:"databases"`
	} `json:"database"`
}

type BackupConfig struct {
	Server_Id        string
	Host             string
	API_Token        string
	Name             string
	Output_Directory string
	Delete_After     int
}

type DatabaseConfig struct {
	Host             string
	User             string
	Password         string
	DBName           string
	Schedule         string
	Output_Directory string
	Delete_After     int
}

func LoadConfig() Config {
	var config Config

	file, err := os.Open("config.json")
	if err != nil {
		logger.ErrorLog.Fatalf("error reading job config file: %v", err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		logger.ErrorLog.Fatalf("error decoding job config: %v", err)
	}

	file.Close()

	return config

}

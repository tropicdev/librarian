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
		} `json:"servers"`
	} `json:"pterodactyl"`
	Database struct {
		Enabled   bool `json:"enabled"`
		Databases []struct {
			Host     string `json:"host"`
			User     string `json:"user"`
			Password string `json:"password"`
			DBName   string `json:"db_name"`
			Schedule string `json:"schedule"`
		} `json:"databases"`
	} `json:"database"`
}

func LoadConfig() Config {
	var config Config

	file, err := os.Open("config.json")
	if err != nil {
		logger.ErrorLog.Fatalf("error reading job config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		logger.ErrorLog.Fatalf("error decoding job config: %v", err)
	}
	return config

}

package backup

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"bitecodelabs.com/librarian/config"
	"bitecodelabs.com/librarian/logger"
)

type Backups struct {
	Object string `json:"object"`
	Data   []struct {
		Object     string `json:"object"`
		Attributes struct {
			UUID         string        `json:"uuid"`
			Name         string        `json:"name"`
			IgnoredFiles []interface{} `json:"ignored_files"`
			Sha256Hash   string        `json:"sha256_hash"`
			Bytes        int           `json:"bytes"`
			CreatedAt    time.Time     `json:"created_at"`
			CompletedAt  time.Time     `json:"completed_at"`
		} `json:"attributes"`
	} `json:"data"`
	Meta struct {
		Pagination struct {
			Total       int `json:"total"`
			Count       int `json:"count"`
			PerPage     int `json:"per_page"`
			CurrentPage int `json:"current_page"`
			TotalPages  int `json:"total_pages"`
			Links       struct {
			} `json:"links"`
		} `json:"pagination"`
	} `json:"meta"`
}

type BackupDetail struct {
	Object     string `json:"object"`
	Attributes struct {
		URL string `json:"url"`
	} `json:"attributes"`
}

func Backup(config config.BackupConfig) {
	backup := getBackups(config)

	backup_urls := getBackupURL(config, backup)

	downloadBackups(backup_urls, config)
}

func getBackups(config config.BackupConfig) Backups {
	var data string
	var backups Backups

	client := &http.Client{}

	req, err := http.NewRequest("GET", config.Host+"/api/client/servers/"+config.Server_Id+"/backups", nil)

	if err != nil {
		logger.ErrorLog.Println(err)
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+config.API_Token)

	resp, err := client.Do(req)

	if err != nil {
		logger.ErrorLog.Println(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	data = string(body)

	if err != nil {
		logger.ErrorLog.Println(err)
	}

	err = json.Unmarshal([]byte(data), &backups)

	if err != nil {
		logger.ErrorLog.Println(err)
	}

	return backups
}

func getBackupURL(config config.BackupConfig, backups Backups) []string {

	var data string
	var backup_urls []string

	client := &http.Client{}

	for _, backup := range backups.Data {
		var backup_data BackupDetail

		req, err := http.NewRequest("GET", config.Host+"/api/client/servers/"+config.Server_Id+"/backups"+backup.Attributes.UUID+"/download", nil)

		if err != nil {
			logger.ErrorLog.Println(err)
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+config.API_Token)

		resp, err := client.Do(req)

		if err != nil {
			logger.ErrorLog.Println(err)
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)

		data = string(body)

		if err != nil {
			logger.ErrorLog.Println(err)
		}

		err = json.Unmarshal([]byte(data), &backup)

		if err != nil {
			logger.ErrorLog.Println(err)
		}

		backup_urls = append(backup_urls, backup_data.Attributes.URL)
	}

	return backup_urls
}

func downloadBackups(backup_urls []string, config config.BackupConfig) {

	for _, backup := range backup_urls {
		client := &http.Client{}

		req, err := http.NewRequest("GET", backup, nil)

		if err != nil {
			logger.ErrorLog.Println(err)
		}

		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+config.API_Token)

		resp, err := client.Do(req)

		if err != nil {
			logger.ErrorLog.Println(err)
		}

		defer resp.Body.Close()

		// Create the file
		out, err := os.Create(config.Output_Directory)

		if err != nil {
			logger.ErrorLog.Println(err)
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		logger.ErrorLog.Println(err)
	}
}

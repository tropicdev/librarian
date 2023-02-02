package backup

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
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

	backup_urls, backup_time := getBackupURL(config, backup)

	downloadBackups(backup_urls, backup_time, config)
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

func getBackupURL(config config.BackupConfig, backups Backups) ([]string, []string) {

	var data string
	var backup_urls []string
	var backup_time []string

	client := &http.Client{}

	for _, backup := range backups.Data {
		var backup_data BackupDetail

		req, err := http.NewRequest("GET", config.Host+"/api/client/servers/"+config.Server_Id+"/backups/"+backup.Attributes.UUID+"/download", nil)

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

		err = json.Unmarshal([]byte(data), &backup_data)

		if err != nil {
			logger.ErrorLog.Println(err)
		}

		backup_time = append(backup_time, backup.Attributes.CreatedAt.Format("02Jan2006_1504"))
		backup_urls = append(backup_urls, backup_data.Attributes.URL)
	}

	return backup_urls, backup_time
}

func downloadBackups(backup_urls []string, backup_time []string, config config.BackupConfig) {

	for i, backup := range backup_urls {
		resp, err := http.Get(backup)

		if err != nil {
			logger.ErrorLog.Println(err)
		}

		defer resp.Body.Close()

		err = os.MkdirAll(config.Output_Directory, 0700)

		if err != nil {
			logger.ErrorLog.Println(err)
		}

		// Create the file
		path := filepath.Join(config.Output_Directory, config.Name+"_"+backup_time[i]+"_"+"to"+"_"+time.Now().Format("02Jan2006_1504")+".tar.gz")
		out, err := os.Create(path)

		if err != nil {
			logger.ErrorLog.Println(err)
		}
		defer out.Close()

		// Write the body to file
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			logger.ErrorLog.Println(err)
		}
	}
}

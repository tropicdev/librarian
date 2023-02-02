package backup

import (
	"time"
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

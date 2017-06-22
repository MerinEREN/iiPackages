package photo

import "time"

// Status is active or deactive
type Photo struct {
	Path         string    `json:"path"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Uploaded     time.Time `json:"uploaded"`
	LastModified time.Time `json:"lastModified"`
	Status       string    `json:"status"`
}

type Photos map[string]*Photo

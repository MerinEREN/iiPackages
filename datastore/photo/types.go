package photo

import "time"

// Photo Status is active or deactive
type Photo struct {
	Link         string    `json:"link"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Uploaded     time.Time `json:"uploaded"`
	LastModified time.Time `json:"lastModified"`
	Status       string    `json:"status"`
}

// Photos is a photos map with encoded photo key as map key.
type Photos map[string]*Photo

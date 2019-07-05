package photo

import "time"

// Photo Status is active or deactive
// "Type" values are main and etc.
type Photo struct {
	ID           string    `datastore:"-"`
	Link         string    `json:"link"`
	Type         string    `json:"type"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	Status       string    `json:"status"`
}

// Photos is a photos map with encoded photo key as map key.
type Photos map[string]*Photo

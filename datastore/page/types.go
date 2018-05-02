package page

import (
	"time"
)

// Page "Exports should have a comment"
// datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type Page struct {
	Title        string    `json:"title"`
	Link         string    `json:"link"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
	ID           string    `datastore:"-"`
}

// Pages "Exports should have a comment"
type Pages map[string]*Page

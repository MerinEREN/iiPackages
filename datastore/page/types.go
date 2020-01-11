package page

import (
	"time"
)

// Page is one of the pages in the app.
// datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type Page struct {
	ID           string    `datastore:"-"`
	Name         string    `json:"name"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}

// Pages "Exports should have a comment"
type Pages map[string]*Page

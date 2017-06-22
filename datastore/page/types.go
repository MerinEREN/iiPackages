package page

import (
	"time"
)

// datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type Page struct {
	ID           string    `datastore:"-"`
	Title        string    `json:"title"`
	Path         string    `json:"path"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}

type Pages map[string]*Page

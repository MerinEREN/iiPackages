package page

import (
	"time"
)

type Page struct {
	ID           string    `datastore:"-"`
	Title        string    `datastore: ",noindex" json:"title"`
	Path         string    `json:"path"`
	Created      time.Time `json:"created"`
	LastModified time.Time `datastore: ",noindex" json:"lastModified"`
}

type Pages []*Page

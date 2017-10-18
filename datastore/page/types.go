package page

import (
	"mime/multipart"
	"time"
)

// Page "Exports should have a comment"
// datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type Page struct {
	ID           string                `datastore:"-"`
	Title        string                `json:"title"`
	Mpf          multipart.File        `datastore:"-" json:"-"`
	Hdr          *multipart.FileHeader `datastore:"-" json:"-"`
	Link         string                `json:"link"`
	Created      time.Time             `json:"created"`
	LastModified time.Time             `json:"lastModified"`
}

// Pages "Exports should have a comment"
type Pages map[string]*Page

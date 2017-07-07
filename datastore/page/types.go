package page

import (
	"mime/multipart"
	"time"
)

// datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type Page struct {
	ID           string                `datastore:"-"`
	Title        string                `json:"title"`
	Mpf          multipart.File        `datastore:"-"`
	Hdr          *multipart.FileHeader `datastore:"-"`
	Link         string                `json:"link"`
	Created      time.Time             `json:"created"`
	LastModified time.Time             `json:"lastModified"`
}

type Pages map[string]*Page

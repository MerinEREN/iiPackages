package language

import (
	"mime/multipart"
	"time"
)

// Language "datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!"
type Language struct {
	ID           string                `datastore:"-"`
	Mpf          multipart.File        `datastore:"-" json:"-"`
	Hdr          *multipart.FileHeader `datastore:"-" json:"-"`
	Link         string                `json:"link"`
	Created      time.Time             `json:"created"`
	LastModified time.Time             `json:"lastModified"`
}

// Languages is a map of *Language with language code as key.
type Languages map[string]*Language

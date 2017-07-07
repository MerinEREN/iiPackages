package language

import (
	"mime/multipart"
	"time"
)

// datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type Language struct {
	Code         string                `datastore:"-" json:"code"`
	Mpf          multipart.File        `datastore:"-"`
	Hdr          *multipart.FileHeader `datastore:"-"`
	Link         string                `json:"link"`
	Created      time.Time             `json:"created"`
	LastModified time.Time             `json:"lastModified"`
}

type Languages map[string]*Language

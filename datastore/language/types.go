package language

import (
	"time"
)

// Language "datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!"
type Language struct {
	ID           string    `datastore:"-"`
	Name         string    `json:"name"`
	Link         string    `json:"link"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}

// Languages is a map of *Language with language code as key.
type Languages map[string]*Language

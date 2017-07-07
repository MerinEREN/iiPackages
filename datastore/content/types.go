package content

import (
	"time"
)

// datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
type Content struct {
	ID           string            `datastore:"-"`
	Values       map[string]string `json:"values"`
	Created      time.Time         `json:"created"`
	LastModified time.Time         `json:"lastModified"`
}

type Contents map[string]*Content

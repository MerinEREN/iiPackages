package content

import (
	"time"
)

// Content is the struct to store application contents for pages with different languages.
// And returns only the requested language as Value.
// datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// Key's "StringID" is content's english value([en-US]). Useful for some cases like
// user's isAdmin() check.
type Content struct {
	ID           string            `datastore:"-"`
	Value        string            `datastore:"-" json:"value"`
	ValuesBS     []byte            `json:"-"`
	Values       map[string]string `datastore:"-" json:"values"`
	Created      time.Time         `json:"created"`
	LastModified time.Time         `json:"lastModified"`
	PageIDs      []string          `datastore:"-" json:"pageIDs"`
}

// Contents is map[string]*Content.
type Contents map[string]*Content

package context

import (
	"time"
)

// Context is the struct to store application contexts for pages with different languages.
// And returns only the requested language as Value.
// datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// Key's "StringID" is a fraction of the context's english value([en-US]).
// Useful for some cases like user's isAdmin() check.
type Context struct {
	ID           string            `datastore:"-"`
	Value        string            `datastore:"-" json:"value"`
	ValuesBS     []byte            `json:"-"`
	Values       map[string]string `datastore:"-" json:"values"`
	Created      time.Time         `json:"created"`
	LastModified time.Time         `json:"lastModified"`
	PageIDs      []string          `datastore:"-" json:"pageIDs"`
}

// Contexts is map[string]*Context.
type Contexts map[string]*Context

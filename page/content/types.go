package content

import (
	"errors"
)

// Errors
var (
	ErrDatastoreDoneContentID = errors.New("Datastore done at PageContent")
	ErrDatastoreDoneContent   = errors.New("Can't find any content for given content " +
		"id.")
)

// en-us, en-au, en-ca, tr...
type Language struct {
	ID     string            `datastore:"-"`
	Values map[string]string `json:"values"`
}

type Languages []Language

type Page struct {
	Name string `json:"name"`
}

type Pages []Page

// Add isInvalidate and set it true for all pages when language changed.
// And add this control at the top of shouldFatch function.
// Group contents by page in Store as root.
// CONSIDERING TO ADD LANG CODE TO THE CONTENT
type Content struct {
	ID     string            `datastore:"-"`
	Values map[string]string `json:"values"`
}

type Contents []Content

// PageID is Parent Key
type PageContent struct {
	ContentKey *datastore.Key
}

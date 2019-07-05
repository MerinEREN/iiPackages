package language

import (
	"time"
)

// Language is an app language.
// "ID" is language code and
// "ContextID" is an encoded "Context" key for multilang usage purpose.
type Language struct {
	ID           string    `datastore:"-"`
	ContextID    string    `json:"contextID"`
	Link         string    `json:"link"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}

// Languages is a map of *Language with language code as key.
type Languages map[string]*Language

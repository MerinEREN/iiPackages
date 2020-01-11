package language

import (
	"time"
)

// Language is an app language.
// "ID" is language code and
// "ContextID" is an encoded "Context" key for multilang usage purpose.
type Language struct {
	ID        string    `datastore:"-"`
	ContextID string    `json:"contextID"`
	Created   time.Time `json:"created"`
}

// Languages is a map of *Language with language code as key.
type Languages map[string]*Language

package tag

import (
	"time"
)

// Tag is an user profession.
// Key's stringID is encoded "Context" key.
// And the "ContextID" is an encoded "Context" key for multilang usage purpose.
type Tag struct {
	ID        string    `datastore:"-"`
	ContextID string    `datastore:"-" json:"contextID"`
	Created   time.Time `json:"created"`
}

// Tags is a map of *Tag with encoded key of Tag as key.
type Tags map[string]*Tag

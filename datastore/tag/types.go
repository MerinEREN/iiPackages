package tag

import (
	"time"
)

// Tag is an account profession.
// And the "Name" is a encoded content key for multilang purpose.
type Tag struct {
	ID      string    `datastore:"-"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
}

// Tags is a map of *Tag with encoded key of Tag as key.
type Tags map[string]*Tag

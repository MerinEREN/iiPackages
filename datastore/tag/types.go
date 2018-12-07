package tag

import (
	"time"
)

// Tag is an user profession.
// Key's stringID is encoded "Content" key.
// And the "ContentID" is an encoded "Content" key for multilang usage purpose.
type Tag struct {
	ID        string    `datastore:"-"`
	ContentID string    `datastore:"-" json:"contentID"`
	Created   time.Time `json:"created"`
}

// Tags is a map of *Tag with encoded key of Tag as key.
type Tags map[string]*Tag

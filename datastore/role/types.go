package role

import (
	"time"
)

// Role represents user roles.
// Key's stringID is encoded "Content" key.
// And the "ContentID" is an encoded "Content" key for multilang usage purpose.
type Role struct {
	ID        string    `datastore:"-"`
	ContentID string    `datastore:"-" json:"contentID"`
	Created   time.Time `json:"-"`
}

// Roles is a map of role pointers with role IDs as their keys.
type Roles map[string]*Role

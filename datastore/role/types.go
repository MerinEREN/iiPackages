package role

import (
	"time"
)

// Role represents user roles.
// Key's stringID is encoded "Context" key.
// And the "ContextID" is an encoded "Context" key for multilang usage purpose.
type Role struct {
	ID        string    `datastore:"-"`
	ContextID string    `datastore:"-" json:"contextID"`
	RoleTypes []string  `datastore:"-" json:"-"`
	Created   time.Time `json:"-"`
}

// Roles is a map of role pointers with role IDs as their keys.
type Roles map[string]*Role

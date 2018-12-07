package roleType

import (
	"time"
)

// RoleType is a type of user's role like "inHouse" and "customer" for now.
// Key's stringID is type's itself.
type RoleType struct {
	ID      string    `datastore:"-"`
	Created time.Time `json:"-"`
}

// RoleTypes is a map of *RoleType with it's ID as key.
type RoleTypes map[string]*RoleType

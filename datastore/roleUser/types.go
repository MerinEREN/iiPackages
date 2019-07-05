package roleUser

import (
	"google.golang.org/appengine/datastore"
)

// RoleUser datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// Encoded role key is key's stringID and user key is the parent key.
type RoleUser struct {
	RoleKey *datastore.Key
}

// RolesUser is a []*RoleUser
type RolesUser []*RoleUser

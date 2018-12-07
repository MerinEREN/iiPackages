package roleTypeRole

import (
	"google.golang.org/appengine/datastore"
)

// RoleTypeRole datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// RoleType key is the parent key.
type RoleTypeRole struct {
	RoleKey *datastore.Key
}

// RoleTypeRoles is a []*RoleTypeRole
type RoleTypeRoles []*RoleTypeRole

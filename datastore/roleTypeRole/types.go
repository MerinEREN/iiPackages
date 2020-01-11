package roleTypeRole

import (
	"google.golang.org/appengine/datastore"
)

// RoleTypeRole datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// "Role" key is the parent key.
type RoleTypeRole struct {
	RoleTypeKey *datastore.Key
}

// RoleTypeRoles is a []*RoleTypeRole
type RoleTypeRoles []*RoleTypeRole

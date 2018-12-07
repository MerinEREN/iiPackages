package userRole

import (
	"google.golang.org/appengine/datastore"
)

// UserRole datastore: ",noindex" causes json naming problems !!!!!!!!!!!!!!!!!!!!!!!!!!!!!
// User key is the parent key.
type UserRole struct {
	RoleKey *datastore.Key
}

// UserRoles is a []*UserRole
type UserRoles []*UserRole

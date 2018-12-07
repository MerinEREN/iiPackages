/*
Package userRole "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
7ne will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package userRole

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetKeysUserOrRole returns the user keys as a slice if the role key is provided
// or returns the role keys as a slice if user key is provided and also an error.
func GetKeysUserOrRole(ctx context.Context, key *datastore.Key) ([]*datastore.Key, error) {
	var kx []*datastore.Key
	q := datastore.NewQuery("UserRole")
	kind := key.Kind()
	switch kind {
	case "Role":
		q = q.
			Filter("RoleKey =", key).
			KeysOnly()
		for it := q.Run(ctx); ; {
			k, err := it.Next(nil)
			if err == datastore.Done {
				return kx, err
			}
			if err != nil {
				return nil, err
			}
			kx = append(kx, k.Parent())
		}
	default:
		// For "User" kind
		q = q.
			Ancestor(key)
		for it := q.Run(ctx); ; {
			ur := new(UserRole)
			_, err := it.Next(ur)
			if err == datastore.Done {
				return kx, err
			}
			if err != nil {
				return nil, err
			}
			kx = append(kx, ur.RoleKey)
		}
	}
}

// Put puts an entity with corresponding entity key and returns an error.
func Put(ctx context.Context, k *datastore.Key, ur *UserRole) error {
	_, err := datastore.Put(ctx, k, ur)
	return err
}

// PutMulti puts entities with corresponding entity keys and returns an error.
func PutMulti(ctx context.Context, kx []*datastore.Key, urx UserRoles) error {
	_, err := datastore.PutMulti(ctx, kx, urx)
	return err
}

// Delete deletes an entity by provided key and returns an error.
func Delete(ctx context.Context, k *datastore.Key) error {
	return datastore.Delete(ctx, k)
}

// GetCount returns the count of the entities that has the provided key and an error.
/* func GetCount(s *session.Session, k *datastore.Key) (c int, err error) {
	q := datastore.NewQuery("UserRole")
	if k.Kind() == "User" {
		q = q.Ancestor(k)
	} else {
		q = q.Filter("RoleKey =", k)
	}
	c, err = q.Count(s.Ctx)
	return
} */

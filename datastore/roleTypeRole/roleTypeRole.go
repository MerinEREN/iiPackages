/*
Package roleTypeRole "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
7ne will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package roleTypeRole

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetKeysRoleTypeOrRole returns the roleType keys as a slice if the role key is provided
// or returns the role keys as a slice if roleType key is provided and also an error.
/* func GetKeysRoleTypeOrRole(ctx context.Context, key *datastore.Key) ([]*datastore.Key, error) {
	var kx []*datastore.Key
	q := datastore.NewQuery("RoleTypeRole")
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
		// For "RoleType" kind
		q = q.
			Ancestor(key)
		for it := q.Run(ctx); ; {
			rtr := new(RoleTypeRole)
			_, err := it.Next(rtr)
			if err == datastore.Done {
				return kx, err
			}
			if err != nil {
				return nil, err
			}
			kx = append(kx, rtr.RoleKey)
		}
	}
} */

/*
GetKeys returns the roleTypeRole keys as a slice if the role key or roleType key
provided, and also an error.
*/
func GetKeys(ctx context.Context, key *datastore.Key) ([]*datastore.Key, error) {
	q := datastore.NewQuery("RoleTypeRole")
	kind := key.Kind()
	switch kind {
	case "RoleType":
		q = q.
			Filter("RoleTypeKey =", key)
	default:
		// For "Role" kind
		q = q.
			Ancestor(key)
	}
	q = q.
		KeysOnly()
	return q.GetAll(ctx, nil)
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

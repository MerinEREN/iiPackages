/*
Package roleUser "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
7ne will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package roleUser

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetKeysByUserOrRoleKey returns the user keys as a slice if the role key is provided
// or returns the role keys as a slice if user key is provided and also an error.
func GetKeysByUserOrRoleKey(ctx context.Context, key *datastore.Key) (
	[]*datastore.Key, error) {
	var kx []*datastore.Key
	q := datastore.NewQuery("RoleUser")
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
			Ancestor(key).
			KeysOnly()
		for it := q.Run(ctx); ; {
			k, err := it.Next(nil)
			if err == datastore.Done {
				return kx, err
			}
			if err != nil {
				return nil, err
			}
			kr, err := datastore.DecodeKey(k.StringID())
			if err != nil {
				return nil, err
			}
			kx = append(kx, kr)
		}
	}
}

// GetKeys returns the roleUser keys by user or role key and an error.
func GetKeys(ctx context.Context, key *datastore.Key) ([]*datastore.Key, error) {
	var kx []*datastore.Key
	q := datastore.NewQuery("RoleUser")
	kind := key.Kind()
	switch kind {
	case "Role":
		q = q.Filter("RoleKey =", key)
	default:
		// For "User" kind
		q = q.Ancestor(key)
	}
	q = q.KeysOnly()
	for it := q.Run(ctx); ; {
		k, err := it.Next(nil)
		if err == datastore.Done {
			return kx, err
		}
		if err != nil {
			return nil, err
		}
		kx = append(kx, k)
	}
}

// GetCount returns the count of the entities that has the provided key and an error.
/* func GetCount(s *session.Session, k *datastore.Key) (c int, err error) {
	q := datastore.NewQuery("RoleUser")
	if k.Kind() == "User" {
		q = q.Ancestor(k)
	} else {
		q = q.Filter("RoleKey =", k)
	}
	c, err = q.Count(s.Ctx)
	return
} */

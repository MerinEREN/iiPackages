/*
Package pageContext "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package pageContext

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetKeysByPageOrContextKey returns the pageContext keys with context keys as a slice
// if the page key is provided or returns the page keys as a slice
// if context key is provided and also an error.
func GetKeysByPageOrContextKey(ctx context.Context, key *datastore.Key) (
	[]*datastore.Key, []*datastore.Key, error) {
	var kx []*datastore.Key
	var kporcx []*datastore.Key
	q := datastore.NewQuery("PageContext")
	kind := key.Kind()
	switch kind {
	case "Context":
		q = q.
			Filter("ContextKey =", key).
			KeysOnly()
		for it := q.Run(ctx); ; {
			k, err := it.Next(nil)
			if err == datastore.Done {
				return kx, kporcx, err
			}
			if err != nil {
				return nil, nil, err
			}
			kx = append(kx, k)
			kporcx = append(kporcx, k.Parent())
		}
	default:
		// "Page" kind
		q = q.
			Ancestor(key)
		for it := q.Run(ctx); ; {
			pc := new(PageContext)
			k, err := it.Next(pc)
			if err == datastore.Done {
				return kx, kporcx, err
			}
			if err != nil {
				return nil, nil, err
			}
			kx = append(kx, k)
			kporcx = append(kporcx, pc.ContextKey)
		}
	}
}

// GetKeysOnly returns corresponding pageContext keys by provided context or page key
// and also returns an error.
func GetKeysOnly(ctx context.Context, k *datastore.Key) ([]*datastore.Key, error) {
	var kx []*datastore.Key
	q := datastore.NewQuery("PageContext")
	if k.Kind() == "Page" {
		q = q.
			Ancestor(k)
	} else {
		q = q.
			Filter("ContextKey =", k)
	}
	q = q.
		KeysOnly()
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

// GetCount returns the count of the entities that has provided key and an error.
func GetCount(ctx context.Context, k *datastore.Key) (c int, err error) {
	q := datastore.NewQuery("PageContext")
	if k.Kind() == "Page" {
		q = q.
			Ancestor(k)
	} else {
		q = q.Filter("ContextKey =", k)
	}
	c, err = q.Count(ctx)
	return
}

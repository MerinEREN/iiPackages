/*
Package pageContent "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package pageContent

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// GetKeysWithPageOrContentKeys returns the pageContent keys with content keys as a slice
// if the page key is provided or returns the page keys as a slice
// if content key is provided and also an error.
func GetKeysWithPageOrContentKeys(ctx context.Context, key *datastore.Key) (
	[]*datastore.Key, []*datastore.Key, error) {
	var kx []*datastore.Key
	var kporcx []*datastore.Key
	q := datastore.NewQuery("PageContent")
	kind := key.Kind()
	switch kind {
	case "Content":
		q = q.
			Filter("ContentKey =", key).
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
			pc := new(PageContent)
			k, err := it.Next(pc)
			if err == datastore.Done {
				return kx, kporcx, err
			}
			if err != nil {
				return nil, nil, err
			}
			kx = append(kx, k)
			kporcx = append(kporcx, pc.ContentKey)
		}
	}
}

// GetKeysOnly returns corresponding pageContent keys by provided content or page key
// and also returns an error.
func GetKeysOnly(ctx context.Context, k *datastore.Key) ([]*datastore.Key, error) {
	var kx []*datastore.Key
	q := datastore.NewQuery("PageContent")
	if k.Kind() == "Page" {
		q = q.
			Ancestor(k)
	} else {
		q = q.
			Filter("ContentKey =", k)
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
	q := datastore.NewQuery("PageContent")
	if k.Kind() == "Page" {
		q = q.
			Ancestor(k)
	} else {
		q = q.Filter("ContentKey =", k)
	}
	c, err = q.Count(ctx)
	return
}

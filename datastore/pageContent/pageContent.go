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
	"github.com/MerinEREN/iiPackages/session"
	"google.golang.org/appengine/datastore"
)

// Get returns the content's keys as a slice if the page ID is provided
// or returns the page's keys as a slice if content ID is provided and also an error.
func Get(s *session.Session, keyEncoded string) ([]*datastore.Key, error) {
	var kx []*datastore.Key
	k, err := datastore.DecodeKey(keyEncoded)
	if err != nil {
		return nil, err
	}
	q := datastore.NewQuery("PageContent")
	if k.Kind() == "Page" {
		q = q.Filter("PageKey =", k).
			Project("ContentKey")
	} else {
		q = q.Filter("ContentKey =", k).
			Project("PageKey")
	}
	for it := q.Run(s.Ctx); ; {
		pc := new(PageContent)
		_, err = it.Next(pc)
		if err == datastore.Done {
			return kx, err
		}
		if err != nil {
			return nil, err
		}
		if k.Kind() == "Page" {
			kx = append(kx, pc.ContentKey)
		} else {
			kx = append(kx, pc.PageKey)
		}
	}
}

// GetKeysOnly returns corresponding pageContent keys by provided content or page key
// and also returns an error.
func GetKeysOnly(s *session.Session, k *datastore.Key) ([]*datastore.Key, error) {
	var pckx []*datastore.Key
	q := datastore.NewQuery("PageContent")
	if k.Kind() == "Page" {
		q = q.Filter("PageKey =", k)
	} else {
		q = q.Filter("ContentKey =", k)
	}
	q = q.KeysOnly()
	for it := q.Run(s.Ctx); ; {
		k, err := it.Next(nil)
		if err == datastore.Done {
			return pckx, err
		}
		if err != nil {
			return nil, err
		}
		pckx = append(pckx, k)
	}
}

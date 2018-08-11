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

// Get returns the content keys as a slice if the page key is provided
// or returns the page keys as a slice if content ID is provided and an error.
// Also as a first argument returns PageContent keys as a slice.
func Get(s *session.Session, key interface{}) ([]*datastore.Key, []*datastore.Key, error) {
	var pckx []*datastore.Key
	var kx []*datastore.Key
	var k *datastore.Key
	var err error
	q := datastore.NewQuery("PageContent")
	switch v := key.(type) {
	case string:
		// Content Kind
		k, err = datastore.DecodeKey(v)
		if err != nil {
			return nil, nil, err
		}
		q = q.Filter("ContentKey =", k).
			Project("PageKey")
	case *datastore.Key:
		// Page Kind
		k = v
		q = q.Filter("PageKey =", k).
			Project("ContentKey")
	}
	for it := q.Run(s.Ctx); ; {
		pc := new(PageContent)
		pck, err := it.Next(pc)
		if err == datastore.Done {
			return pckx, kx, err
		}
		if err != nil {
			return nil, nil, err
		}
		if k.Kind() == "Page" {
			kx = append(kx, pc.ContentKey)
		} else {
			kx = append(kx, pc.PageKey)
		}
		pckx = append(pckx, pck)
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

// GetCount returns the count of the entities that has provided key and an error.
func GetCount(s *session.Session, k *datastore.Key) (c int, err error) {
	q := datastore.NewQuery("PageContent")
	if k.Kind() == "Page" {
		q = q.Filter("PageKey =", k)
	} else {
		q = q.Filter("ContentKey =", k)
	}
	c, err = q.Count(s.Ctx)
	return
}

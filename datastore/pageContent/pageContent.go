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

// Get returns the content's key's intID as a slice of string if pID is provided
// or returns PageIDs as a slice if cID is provided and also an error.
func Get(s *session.Session, pID, cID string) ([]string, error) {
	var IDx []string
	var err error
	q := datastore.NewQuery("PageContent")
	if pID != "" {
		q = q.Filter("PageID =", pID).
			Project("ContentID")
	} else if cID != "" {
		q = q.Filter("ContentID =", cID).
			Project("PageID")
	}
	for it := q.Run(s.Ctx); ; {
		pc := new(PageContent)
		_, err = it.Next(pc)
		if err == datastore.Done {
			return IDx, err
		}
		if err != nil {
			return nil, err
		}
		if pID != "" {
			IDx = append(IDx, pc.ContentID)
		} else if cID != "" {
			IDx = append(IDx, pc.PageID)
		}
	}
}

// DeleteMulti deletes all entities for provided content keys.
func DeleteMulti(s *session.Session, kx []*datastore.Key) error {
	it := new(datastore.Iterator)
	k := new(datastore.Key)
	var err error
	q := datastore.NewQuery("PageContent")
	for _, v := range kx {
		q.Filter("ContentID =", v).
			KeysOnly()
		for it = q.Run(s.Ctx); ; {
			k, err = it.Next(nil)
			if err == datastore.Done {
				break
			}
			if err != nil {
				return err
			}
			err = datastore.Delete(s.Ctx, k)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

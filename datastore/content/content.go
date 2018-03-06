/*
Package content "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package content

import (
	"github.com/MerinEREN/iiPackages/datastore/pageContent"
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"strconv"
	"time"
)

/*
GetMulti returns corresponding page contents with all languages if content keys provided
Otherwise returns limited entitity from the given cursor.
If limit is nil default limit will be used.
*/
func GetMulti(s *session.Session, crsr datastore.Cursor, limit, kx interface{}) (
	Contents, datastore.Cursor, error) {
	cs := make(Contents)
	var err error
	var cx []*Content
	if kx, ok := kx.([]*datastore.Key); ok {
		// RETURNED ENTITY LIMIT COUD BE A PROBLEM HERE !!!!!!!!!!!!!!!!!!!!!!!!!!!
		if err = datastore.GetMulti(s.Ctx, kx, cx); err != nil {
			return nil, crsr, err
		}
		for i, v := range kx {
			cs[strconv.FormatInt(v.IntID(), 10)] = cx[i]
		}
		return cs, crsr, err
	}
	q := datastore.NewQuery("Content").
		Order("-LastModified")
	if crsr.String() != "" {
		q = q.Start(crsr)
	}
	if limit != nil {
		l := limit.(int)
		q = q.Limit(l)
	} else {
		q = q.Limit(20)
	}
	for it := q.Run(s.Ctx); ; {
		c := new(Content)
		k, err := it.Next(c)
		if err == datastore.Done {
			crsr, err = it.Cursor()
			return cs, crsr, err
		}
		if err != nil {
			return nil, crsr, err
		}
		c.ID = strconv.FormatInt(k.IntID(), 10)
		cs[c.ID] = c
	}
}

/*
PutMulti is a transaction that delets all PageContent entities for corresponding content
if only the request method is "PUT".
And then creates and puts new ones.
Finally, puts modified or newly created Contents and returns new Contents if request method
is POST..
*/
func PutMulti(s *session.Session, cx []*Content) (Contents, error) {
	var kx []*datastore.Key
	var keyIntID int64
	var err error
	for _, v := range cx {
		k := new(datastore.Key)
		if s.R.Method == "PUT" {
			keyIntID, err = strconv.ParseInt(v.ID, 10, 64)
			if err != nil {
				return nil, err
			}
			k = datastore.NewKey(s.Ctx, "Content", "", keyIntID, nil)
		} else {
			k = datastore.NewIncompleteKey(s.Ctx, "Content", nil)
			v.Created = time.Now()
		}
		kx = append(kx, k)
		v.LastModified = time.Now()
	}
	err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		if s.R.Method == "PUT" {
			err1 = pageContent.DeleteMulti(s, kx)
			if err1 != nil {
				return err1
			}
		}
		// Can not store incomplete keys in datastore.
		// ALWAYS USE RETURNED KEYS.
		kx, err1 = datastore.PutMulti(s.Ctx, kx, cx)
		if err1 != nil {
			return err1
		}
		for i, v := range cx {
			for _, v2 := range v.Pages {
				pcK := datastore.NewIncompleteKey(s.Ctx, "PageContent", nil)
				pc := new(pageContent.PageContent)
				pc.ContentID = strconv.FormatInt(kx[i].IntID(), 10)
				pc.PageID = v2
				_, err1 = datastore.Put(s.Ctx, pcK, pc)
				if err1 != nil {
					return err1
				}
			}
		}
		return err1
	}, nil)
	if err != nil {
		return nil, err
	}
	if s.R.Method == "POST" {
		cs := make(Contents)
		for i, v := range cx {
			v.ID = strconv.FormatInt(kx[i].IntID(), 10)
			cs[v.ID] = v
		}
		return cs, err
	}
	return nil, err
}

// PutMultiAndGetMulti is a transaction which puts the posted entities first
// and then gets entities from the reseted cursor with the given limit.
// Finally returnes received entities with posted entities added to them
// as a map.
func PutMultiAndGetMulti(s *session.Session, c datastore.Cursor, cx []*Content) (
	Contents, datastore.Cursor, error) {
	csPut := make(Contents)
	csGet := make(Contents)
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		if csPut, err1 = PutMulti(s, cx); err1 != nil {
			return err1
		}
		if csGet, c, err1 = GetMulti(s, c, nil, nil); err1 == nil {
			for i, v := range csGet {
				csPut[i] = v
			}
		}
		return err1
	}, nil)
	return csPut, c, err
}

// Delete removes the entity with the provided IntID as string and returns an error..
func Delete(s *session.Session, keyIntID string) error {
	IntID, err := strconv.ParseInt(keyIntID, 10, 64)
	if err != nil {
		return err
	}
	k := datastore.NewKey(s.Ctx, "Content", "", IntID, nil)
	return datastore.Delete(s.Ctx, k)
}

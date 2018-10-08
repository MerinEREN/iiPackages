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
	"time"
)

/*
GetMulti returns corresponding entities with values of all languages
if keys provided.
Otherwise returns limited entitities from the given cursor.
If limit is nil default limit will be used.
*/
func GetMulti(s *session.Session, crsr datastore.Cursor, limit, kx interface{}) (
	Contents, datastore.Cursor, error) {
	cs := make(Contents)
	if kx, ok := kx.([]*datastore.Key); ok {
		cx := make([]*Content, len(kx))
		// RETURNED ENTITY LIMIT COULD BE A PROBLEM HERE !!!!!!!!!!!!!!!!!!!!!!!!!!
		err := datastore.GetMulti(s.Ctx, kx, cx)
		if err != nil {
			return nil, crsr, err
		}
		for i, v := range kx {
			cs[v.Encode()] = cx[i]
		}
		return cs, crsr, err
	}
	q := datastore.NewQuery("Content").
		Order("-LastModified")
	if crsr.String() != "" {
		q = q.Start(crsr)
	}
	if limit != nil {
		// l := limit.(int)
		// q = q.Limit(l)
	} else {
		// q = q.Limit(20)
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
		c.ID = k.Encode()
		cs[c.ID] = c
	}
}

/*
PutMulti is a transaction that delets all PageContent entities for corresponding content
if only the request method is "PUT".
And then creates and puts new ones.
Finally, puts modified or newly created Contents and returns new Contents if request method
is POST.
*/
func PutMulti(s *session.Session, cx []*Content) (Contents, error) {
	var kx []*datastore.Key
	var pckx []*datastore.Key
	var err error
	for _, v := range cx {
		k := new(datastore.Key)
		if s.R.Method == "PUT" {
			k, err = datastore.DecodeKey(v.ID)
			if err != nil {
				return nil, err
			}
			pckx2, err := pageContent.GetKeysOnly(s, k)
			if err != datastore.Done {
				return nil, err
			}
			for _, v2 := range pckx2 {
				pckx = append(pckx, v2)
			}
		} else {
			k = datastore.NewIncompleteKey(s.Ctx, "Content", nil)
			v.Created = time.Now()
		}
		kx = append(kx, k)
		v.LastModified = time.Now()
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		if s.R.Method == "PUT" {
			err1 = datastore.DeleteMulti(ctx, pckx)
			if err1 != nil {
				return
			}
		}
		// Can not store incomplete keys in datastore.
		// ALWAYS USE RETURNED KEYS.
		kx, err1 = datastore.PutMulti(ctx, kx, cx)
		if err1 != nil {
			return
		}
		for i, v := range cx {
			for _, v2 := range v.PageIDs {
				pck := datastore.NewIncompleteKey(ctx, "PageContent", nil)
				pc := new(pageContent.PageContent)
				pc.ContentKey = kx[i]
				pk := new(datastore.Key)
				pk, err1 = datastore.DecodeKey(v2)
				if err1 != nil {
					return
				}
				pc.PageKey = pk
				_, err1 = datastore.Put(ctx, pck, pc)
				if err1 != nil {
					return
				}
			}
		}
		return
	}, opts)
	if err != nil {
		return nil, err
	}
	if s.R.Method == "POST" {
		cs := make(Contents)
		for i, v := range cx {
			v.ID = kx[i].Encode()
			cs[v.ID] = v
		}
		return cs, err
	}
	return nil, err
}

// PutMultiAndGetMulti is a transaction which puts the posted entities first
// and then gets entities from the reseted cursor by the given limit.
// Finally returnes received entities with posted entities added to them
// as a map.
func PutMultiAndGetMulti(s *session.Session, crsr datastore.Cursor, cx []*Content) (
	Contents, datastore.Cursor, error) {
	csPut := make(Contents)
	csGet := make(Contents)
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		if csPut, err1 = PutMulti(s, cx); err1 != nil {
			return
		}
		if csGet, crsr, err1 = GetMulti(s, crsr, nil, nil); err1 == nil ||
			err1 == datastore.Done {
			for i, v := range csGet {
				csPut[i] = v
			}
		}
		return
	}, nil)
	return csPut, crsr, err
}

// DeleteMulti removes the entities
// and all the corresponding pageContent entities by the provided encoded keys
// also returns an error.
func DeleteMulti(s *session.Session, ekx []string) error {
	var kx []*datastore.Key
	var pckx []*datastore.Key
	for _, v := range ekx {
		k, err := datastore.DecodeKey(v)
		if err != nil {
			return err
		}
		kx = append(kx, k)
		pckx2, err := pageContent.GetKeysOnly(s, k)
		if err != datastore.Done {
			return err
		}
		for _, v2 := range pckx2 {
			pckx = append(pckx, v2)
		}
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	return datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		err1 = datastore.DeleteMulti(ctx, pckx)
		if err1 != nil {
			return
		}
		err1 = datastore.DeleteMulti(ctx, kx)
		return
	}, opts)
}

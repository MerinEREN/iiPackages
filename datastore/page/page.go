/*
Package page "Every package should have a package comment, a block comment preceding
the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package page

import (
	"errors"
	"github.com/MerinEREN/iiPackages/datastore/pageContent"
	"github.com/MerinEREN/iiPackages/session"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"strings"
	"time"
)

// Errors
var (
	ErrFindPage = errors.New("error while getting page")
)

/*
GetMulti returns corresponding entities if keys provided
Otherwise returns limited entitities from the given cursor.
If limit is nil default limit will be used.
*/
func GetMulti(s *session.Session, crsr datastore.Cursor, limit, kx interface{}) (
	Pages, datastore.Cursor, error) {
	ps := make(Pages)
	if kx, ok := kx.([]*datastore.Key); ok {
		px := make([]*Page, len(kx))
		// RETURNED ENTITY LIMIT COULD BE A PROBLEM HERE !!!!!!!!!!!!!!!!!!!!!!!!!!
		err := datastore.GetMulti(s.Ctx, kx, px)
		if err != nil {
			return nil, crsr, err
		}
		for i, v := range kx {
			ps[v.Encode()] = px[i]
		}
		return ps, crsr, err
	}
	// Maybe -LastModied should be the order criteria if consider UX, think about that.
	q := datastore.NewQuery("Page").
		Project("Text", "Link").
		Order("-Created")
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
		p := new(Page)
		k, err := it.Next(p)
		if err == datastore.Done {
			crsr, err = it.Cursor()
			return ps, crsr, err
		}
		if err != nil {
			err = ErrFindPage
			return ps, crsr, err
		}
		p.ID = k.Encode()
		ps[p.ID] = p
	}
}

// Put adds to or modifies an entity in the kind according to request method.
func Put(s *session.Session, p *Page) (*Page, error) {
	k := new(datastore.Key)
	var err error
	p.LastModified = time.Now()
	if s.R.Method == "POST" {
		stringID := strings.ToLower(strings.Replace(p.Text, " ", "", -1))
		k = datastore.NewKey(s.Ctx, "Page", stringID, 0, nil)
		p.Created = time.Now()
		k, err = datastore.Put(s.Ctx, k, p)
		p.ID = k.Encode()
	} else if s.R.Method == "PUT" {
		k, err = datastore.DecodeKey(p.ID)
		if err != nil {
			return nil, err
		}
		tempP := new(Page)
		err = datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (
			err1 error) {
			if err1 = datastore.Get(ctx, k, tempP); err1 != nil {
				return
			}
			p.Created = tempP.Created
			_, err1 = datastore.Put(ctx, k, p)
			return
		}, nil)
	}
	return p, err
}

// PutAndGetMulti is a transaction which puts the posted item first
// and then gets entities with the given limit.
func PutAndGetMulti(s *session.Session, crsr datastore.Cursor, p *Page) (Pages,
	datastore.Cursor, error) {
	ps := make(Pages)
	pNew := new(Page)
	// USAGE "s" INSTEAD OF "ctx" INSIDE THE TRANSACTION IS WRONG !!!!!!!!!!!!!!!!!!!!!
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		pNew, err1 = Put(s, p)
		if err1 != nil {
			return
		}
		ps, crsr, err1 = GetMulti(s, crsr, 9, nil)
		return
	}, nil)
	ps[pNew.ID] = pNew
	return ps, crsr, err
}

// Get returns the page and an error by given encoded key.
func Get(s *session.Session, keyEncoded string) (Pages, error) {
	ps := make(Pages)
	k, err := datastore.DecodeKey(keyEncoded)
	if err != nil {
		return nil, err
	}
	p := new(Page)
	err = datastore.Get(s.Ctx, k, p)
	p.ID = keyEncoded
	ps[p.ID] = p
	return ps, err
}

// Delete removes the entity and all the corresponding pageContent entities
// by the provided encoded key
// and if a content only been included in that page also gonna be removed.
// As a return returns an error.
func Delete(s *session.Session, keyEncoded string) error {
	k, err := datastore.DecodeKey(keyEncoded)
	if err != nil {
		return err
	}
	kpcx, kcx, err := pageContent.GetKeysWithPageOrContentKeys(s.Ctx, k)
	if err != datastore.Done {
		return err
	}
	var count int
	var kcx2 []*datastore.Key
	for _, v := range kcx {
		count, err = pageContent.GetCount(s.Ctx, v)
		if err != nil {
			return err
		}
		if count == 1 {
			kcx2 = append(kcx2, v)
		}
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	return datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		err1 = datastore.DeleteMulti(ctx, kpcx)
		if err1 != nil {
			return
		}
		err1 = datastore.DeleteMulti(ctx, kcx2)
		if err1 != nil {
			return
		}
		err1 = datastore.Delete(ctx, k)
		return
	}, opts)
}

// DeleteMulti removes the entities and all the corresponding pageContent entities
// by the provided encoded keys.
// And if a content only been included in one of the deleted pages also gonna be removed.
// As a return returns an error.
func DeleteMulti(s *session.Session, ekx []string) error {
	var kx []*datastore.Key
	var kpcx []*datastore.Key
	var kcx []*datastore.Key
	var count int
	var kpcx2 []*datastore.Key
	var kcx2 []*datastore.Key
	for _, v := range ekx {
		k, err := datastore.DecodeKey(v)
		if err != nil {
			return err
		}
		kx = append(kx, k)
		kpcx2, kcx2, err = pageContent.GetKeysWithPageOrContentKeys(s.Ctx, k)
		if err != datastore.Done {
			return err
		}
		for _, v2 := range kpcx2 {
			kpcx = append(kpcx, v2)
		}
		for _, v3 := range kcx2 {
			count, err = pageContent.GetCount(s.Ctx, v3)
			if err != nil {
				return err
			}
			if count == 1 {
				kcx = append(kcx, v3)
			}
		}
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	return datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		err1 = datastore.DeleteMulti(ctx, kpcx)
		if err1 != nil {
			return
		}
		err1 = datastore.DeleteMulti(ctx, kcx)
		if err1 != nil {
			return
		}
		err1 = datastore.DeleteMulti(ctx, kx)
		return
	}, opts)
}

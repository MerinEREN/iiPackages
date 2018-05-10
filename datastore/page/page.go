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
	"time"
)

// Errors
var (
	ErrFindPage = errors.New("error while getting page")
)

/*
GetMulti returns corresponding pages if page keys provided
Otherwise returns limited entitity from the given cursor.
If limit is nil default limit will be used.
*/
func GetMulti(s *session.Session, c datastore.Cursor, limit, kx interface{}) (Pages, datastore.Cursor, error) {
	ps := make(Pages)
	if kx, ok := kx.([]*datastore.Key); ok {
		var px []*Page
		err := datastore.GetMulti(s.Ctx, kx, px)
		if err != nil {
			return nil, c, err
		}
		for i, v := range kx {
			px[i].ID = v.Encode()
			ps[px[i].ID] = px[i]
		}
		return ps, c, err
	}
	// Maybe -LastModied should be the order ctireia if consider UX, think about that.
	q := datastore.NewQuery("Page").
		Project("Title", "Link").
		Order("-Created")
	if c.String() != "" {
		q = q.Start(c)
	}
	if limit != nil {
		l := limit.(int)
		q = q.Limit(l)
	} else {
		q = q.Limit(10)
	}
	for it := q.Run(s.Ctx); ; {
		p := new(Page)
		k, err := it.Next(p)
		if err == datastore.Done {
			c, err = it.Cursor()
			return ps, c, err
		}
		if err != nil {
			err = ErrFindPage
			return nil, c, err
		}
		p.ID = k.Encode()
		ps[p.ID] = p
	}
}

/*
Put "Inside a package, any comment immediately preceding a top-level declaration serves as a
doc comment for that declaration. Every exported (capitalized) name in a program should
have a doc comment.
Doc comments work best as complete sentences, which allow a wide variety of automated
presentations. The first sentence should be a one-sentence summary that starts with the
name being declared."
*/
func Put(s *session.Session, p *Page) (*Page, error) {
	k := new(datastore.Key)
	var err error
	p.LastModified = time.Now()
	if s.R.Method == "POST" {
		k = datastore.NewIncompleteKey(s.Ctx, "Page", nil)
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
func PutAndGetMulti(s *session.Session, c datastore.Cursor, p *Page) (Pages,
	datastore.Cursor, error) {
	ps := make(Pages)
	pNew := new(Page)
	err := datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		pNew, err1 = Put(s, p)
		if err1 != nil {
			return
		}
		ps, c, err1 = GetMulti(s, c, 9, nil)
		return
	}, nil)
	ps[pNew.ID] = pNew
	return ps, c, err
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
	pckx, ckx, err := pageContent.Get(s, keyEncoded)
	if err != datastore.Done {
		return err
	}
	var count int
	var ckx2 []*datastore.Key
	for _, v := range ckx {
		count, err = pageContent.GetCount(s, v)
		if err != nil {
			return err
		}
		if count == 1 {
			ckx2 = append(ckx2, v)
		}
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	return datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		err1 = datastore.DeleteMulti(ctx, pckx)
		if err1 != nil {
			return
		}
		err1 = datastore.DeleteMulti(ctx, ckx2)
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
	var pckx []*datastore.Key
	var ckx []*datastore.Key
	var count int
	var pckx2 []*datastore.Key
	var ckx2 []*datastore.Key
	for _, v := range ekx {
		k, err := datastore.DecodeKey(v)
		if err != nil {
			return err
		}
		kx = append(kx, k)
		pckx2, ckx2, err = pageContent.Get(s, v)
		if err != datastore.Done {
			return err
		}
		for _, v2 := range pckx2 {
			pckx = append(pckx, v2)
		}
		for _, v3 := range ckx2 {
			count, err = pageContent.GetCount(s, v3)
			if err != nil {
				return err
			}
			if count == 1 {
				ckx = append(ckx, v3)
			}
		}
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	return datastore.RunInTransaction(s.Ctx, func(ctx context.Context) (err1 error) {
		err1 = datastore.DeleteMulti(ctx, pckx)
		if err1 != nil {
			return
		}
		err1 = datastore.DeleteMulti(ctx, ckx)
		if err1 != nil {
			return
		}
		err1 = datastore.DeleteMulti(ctx, kx)
		return
	}, opts)
}

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
	"github.com/MerinEREN/iiPackages/datastore/pageContext"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"time"
)

// Errors
var (
	ErrFindPage = errors.New("error while getting page")
)

// GetAll returns all the entities with some projections from the begining of the kind.
func GetAll(ctx context.Context) (Pages, error) {
	var px []*Page
	q := datastore.NewQuery("Page")
	q = q.
		Project("Name").
		Order("-Created")
	kx, err := q.GetAll(ctx, px)
	if err != nil {
		return nil, err
	}
	ps := make(Pages)
	for i, v := range kx {
		px[i].ID = v.Encode()
		ps[v.Encode()] = px[i]
	}
	return ps, err
}

/*
Update updates and returns (only with "ID" and "LastModified") the entity
by the given encoded entity key "ek" and also returns an error.
*/
func Update(ctx context.Context, p *Page, ek string) (*Page, error) {
	k, err := datastore.DecodeKey(ek)
	if err != nil {
		return nil, err
	}
	p2 := new(Page)
	if err = datastore.Get(ctx, k, p2); err != nil {
		return nil, err
	}
	p.Created = p2.Created
	p.LastModified = time.Now()
	_, err = datastore.Put(ctx, k, p)
	p3 := &Page{
		ID:           ek,
		LastModified: p.LastModified,
	}
	return p3, err
}

/*
Delete removes the entity and all the corresponding pageContext entities
by the provided encoded key
and if a context only been included in that page also gonna be removed.
As a return returns an error.
*/
func Delete(ctx context.Context, ek string) error {
	k, err := datastore.DecodeKey(ek)
	if err != nil {
		return err
	}
	kpcx, kcx, err := pageContext.GetKeysByPageOrContextKey(ctx, k)
	if err != datastore.Done {
		return err
	}
	var count int
	var kcx2 []*datastore.Key
	for _, v := range kcx {
		count, err = pageContext.GetCount(ctx, v)
		if err != nil {
			return err
		}
		if count == 1 {
			kcx2 = append(kcx2, v)
		}
	}
	opts := new(datastore.TransactionOptions)
	opts.XG = true
	return datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
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

/*
DeleteMulti removes the entities and all the corresponding pageContext entities
by the provided encoded keys.
And if a context only been included in one of the deleted pages also gonna be removed.
As a return returns an error.
*/
func DeleteMulti(ctx context.Context, ekx []string) error {
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
		kpcx2, kcx2, err = pageContext.GetKeysByPageOrContextKey(ctx, k)
		if err != datastore.Done {
			return err
		}
		for _, v2 := range kpcx2 {
			kpcx = append(kpcx, v2)
		}
		for _, v3 := range kcx2 {
			count, err = pageContext.GetCount(ctx, v3)
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
	return datastore.RunInTransaction(ctx, func(ctx context.Context) (err1 error) {
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

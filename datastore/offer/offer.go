/*
Package offer "Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows."
*/
package offer

import (
	dstore "github.com/MerinEREN/iiPackages/datastore"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
)

// Get returns limited entities from the start of the kind with given filters and order.
// Add location filter here and use accounts addres info.
func Get(ctx context.Context, uTagIDs []*datastore.Key) (
	Offers, datastore.Cursor, datastore.Cursor, error) {
	os := make(Offers)
	var cEnd datastore.Cursor
	q := datastore.NewQuery("Offer")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		Limit(10)
	it := q.Run(ctx)
	cStart, err := it.Cursor()
	if err != nil {
		log.Printf("Offer Get Cursor Error: %v\n", err)
	}
	for {
		o := new(Offer)
		k, err := it.Next(o)
		if err == datastore.Done {
			cEnd, err = it.Cursor()
			return os, cStart, cEnd, err
		}
		if err != nil {
			return nil, cStart, cEnd, err
		}
		o.ID = k.StringID()
		os[o.ID] = o
	}
}

// GetNewest returns all the results from the begining to the previous start point
// with given filters and order.
// Add location filter here and use accounts addres info.
func GetNewest(ctx context.Context, cStart datastore.Cursor, uTagIDs []*datastore.Key) (
	Offers, datastore.Cursor, error) {
	os := make(Offers)
	q := datastore.NewQuery("Offer")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		End(cStart)
	it := q.Run(ctx)
	cStart, err := it.Cursor()
	if err != nil {
		log.Printf("Offer GetNewest Cursor Error: %v\n", err)
	}
	for {
		o := new(Offer)
		k, err := it.Next(o)
		if err == datastore.Done {
			return os, cStart, err
		}
		if err != nil {
			return nil, cStart, err
		}
		o.ID = k.StringID()
		os[o.ID] = o
	}
}

// GetNewestCount returns the results count from the begining to the previous start point
// with given filters and order.
// Add location filter here and use accounts addres info.
func GetNewestCount(ctx context.Context, cStart datastore.Cursor, uTagIDs []*datastore.Key) (
	cnt int, err error) {
	q := datastore.NewQuery("Offer")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		End(cStart)
	cnt, err = q.Count(ctx)
	return
}

// GetOldest returns limited entities from the previous cEnd cursor
// with given filters and order.
// Add location filter here and use accounts addres info.
func GetOldest(ctx context.Context, cEnd datastore.Cursor, uTagIDs []*datastore.Key) (
	Offers, datastore.Cursor, error) {
	os := make(Offers)
	q := datastore.NewQuery("Offer")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		Start(cEnd).
		Limit(10)
	for it := q.Run(ctx); ; {
		o := new(Offer)
		k, err := it.Next(o)
		if err == datastore.Done {
			cEnd, err = it.Cursor()
			return os, cEnd, err
		}
		if err != nil {
			return nil, cEnd, err
		}
		o.ID = k.StringID()
		os[o.ID] = o
	}
}

// GetViaParent returns limited entities from the previous cursor
// with given filters and order.
func GetViaParent(ctx context.Context, crsr datastore.Cursor, pk *datastore.Key) (
	Offers, datastore.Cursor, error) {
	os := make(Offers)
	q := datastore.NewQuery("Offer")
	q = q.Ancestor(pk).
		Order("-LastModified").
		Start(crsr).
		Limit(10)
	for it := q.Run(ctx); ; {
		o := new(Offer)
		k, err := it.Next(o)
		if err == datastore.Done {
			crsr, err = it.Cursor()
			return os, crsr, err
		}
		if err != nil {
			return nil, crsr, err
		}
		o.ID = k.Encode()
		os[o.ID] = o
	}
}

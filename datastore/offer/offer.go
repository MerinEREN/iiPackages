/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package offer

import (
	dstore "github.com/MerinEREN/iiPackages/datastore"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
)

// Add location filter here and use accounts addres info.
func Get(ctx context.Context, uTagIDs []*datastore.Key) (
	Offers, datastore.Cursor, datastore.Cursor, error) {
	o := new(Offer)
	var os Offers
	var cOld datastore.Cursor
	q := datastore.NewQuery("Offer")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		Limit(10)
	it := q.Run(ctx)
	cNew, err := it.Cursor()
	if err != nil {
		log.Printf("Offer Get Error: %v\n", err)
	}
	for {
		k, err := it.Next(o)
		if err == datastore.Done {
			cOld, err = it.Cursor()
			return os, cNew, cOld, err
		}
		if err != nil {
			return nil, cNew, cOld, err
		}
		o.ID = k.StringID()
		os = append(os, o)
	}
}

// Add location filter here and use accounts addres info.
func GetNewest(ctx context.Context, c datastore.Cursor, uTagIDs []*datastore.Key) (
	Offers, datastore.Cursor, error) {
	o := new(Offer)
	var os Offers
	q := datastore.NewQuery("Offer")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("LastModified").
		Start(c)
	for it := q.Run(ctx); ; {
		k, err := it.Next(o)
		if err == datastore.Done {
			c, err = it.Cursor()
			return os, c, err
		}
		if err != nil {
			return nil, c, err
		}
		o.ID = k.StringID()
		os = append(os, o)
	}
}

// Add location filter here and use accounts addres info.
func GetNewestCount(ctx context.Context, c datastore.Cursor, uTagIDs []*datastore.Key) (
	cnt int, err error) {
	q := datastore.NewQuery("Offer")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("LastModified").
		Start(c)
	cnt, err = q.Count(ctx)
	return
}

// Add location filter here and use accounts addres info.
func GetOldest(ctx context.Context, c datastore.Cursor, uTagIDs []*datastore.Key) (
	Offers, datastore.Cursor, error) {
	o := new(Offer)
	var os Offers
	q := datastore.NewQuery("Offer")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		Start(c).
		Limit(10)
	for it := q.Run(ctx); ; {
		k, err := it.Next(o)
		if err == datastore.Done {
			c, err = it.Cursor()
			return os, c, err
		}
		if err != nil {
			return nil, c, err
		}
		o.ID = k.StringID()
		os = append(os, o)
	}
}

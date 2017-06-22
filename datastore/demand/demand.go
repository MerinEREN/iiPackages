/*
Every package should have a package comment, a block comment preceding the package clause.
For multi-file packages, the package comment only needs to be present in one file, and any
one will do. The package comment should introduce the package and provide information
relevant to the package as a whole. It will appear first on the godoc page and should set
up the detailed documentation that follows.
*/
package demand

import (
	dstore "github.com/MerinEREN/iiPackages/datastore"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"log"
)

// Add location filter here and use accounts addres info.
func Get(ctx context.Context, uTagIDs []*datastore.Key) (
	Demands, datastore.Cursor, datastore.Cursor, error) {
	ds := make(Demands)
	var cOld datastore.Cursor
	q := datastore.NewQuery("Demand")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		Limit(10)
	it := q.Run(ctx)
	cNew, err := it.Cursor()
	if err != nil {
		log.Printf("Demand Get Error: %v\n", err)
	}
	for {
		d := new(Demand)
		k, err := it.Next(d)
		if err == datastore.Done {
			cOld, err = it.Cursor()
			return ds, cNew, cOld, err
		}
		if err != nil {
			return nil, cNew, cOld, err
		}
		d.ID = k.StringID()
		ds[d.ID] = d
	}
}

// Add location filter here and use accounts addres info.
func GetNewest(ctx context.Context, c datastore.Cursor, uTagIDs []*datastore.Key) (
	Demands, datastore.Cursor, error) {
	ds := make(Demands)
	q := datastore.NewQuery("Demand")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("LastModified").
		Start(c)
	for it := q.Run(ctx); ; {
		d := new(Demand)
		k, err := it.Next(d)
		if err == datastore.Done {
			c, err = it.Cursor()
			return ds, c, err
		}
		if err != nil {
			return nil, c, err
		}
		d.ID = k.StringID()
		ds[d.ID] = d
	}
}

// Add location filter here and use accounts addres info.
func GetNewestCount(ctx context.Context, c datastore.Cursor, uTagIDs []*datastore.Key) (
	cnt int, err error) {
	q := datastore.NewQuery("Demand")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("LastModified").
		Start(c)
	cnt, err = q.Count(ctx)
	return
}

// Add location filter here and use accounts addres info.
func GetOldest(ctx context.Context, c datastore.Cursor, uTagIDs []*datastore.Key) (
	Demands, datastore.Cursor, error) {
	ds := make(Demands)
	q := datastore.NewQuery("Demand")
	q = dstore.FilterMulti(q, "TagIDs =", uTagIDs).
		Order("-LastModified").
		Start(c).
		Limit(10)
	for it := q.Run(ctx); ; {
		d := new(Demand)
		k, err := it.Next(d)
		if err == datastore.Done {
			c, err = it.Cursor()
			return ds, c, err
		}
		if err != nil {
			return nil, c, err
		}
		d.ID = k.StringID()
		ds[d.ID] = d
	}
}
